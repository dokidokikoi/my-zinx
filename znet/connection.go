package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/dokidokikoi/my-zinx/ziface"
)

type Connection struct {
	// 当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 当前连接的 ID, 也可以称为 SessionID, ID 全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool

	// 消息管理 MsgID 和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandler

	// 告知该连接已经退出/停止的 channel
	ExitBuffChan chan bool
}

// 停止连接，结束当前连接状态 M
func (c *Connection) Stop() {
	// 1. 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// TODO: Connection.Stop() 如果用户注册了该连接的关闭回调业务，
	// 则在此时应该显示调用

	// 关闭 Socket 连接
	c.Conn.Close()

	// 通知从缓冲队列读数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	// 关闭该连接的全部的管道
	close(c.ExitBuffChan)
}

// 处理 conn 读数据的 Gorutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goruntine is running")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		// 创建拆包对象
		dp := NewDataPack()

		// 读取客户端的 Msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		// 拆包
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}

		// 根据 dataLen 读取数据
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		// 从绑定好的消息和对应的处理方法中执行对应的 Handle 方法
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	// 开启处理该连接读取客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息,不再阻塞
			return
		}
	}
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	// 将 data 封包，并发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Println("pack error msg id = ", msgID)
		return errors.New("Pack error msg")
	}

	// 写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id ", msgID, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}

	return nil
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   msgHandler,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}
