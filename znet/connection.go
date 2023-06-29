package znet

import (
	"fmt"
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

	// 该连接的处理方式 API
	handleAPI ziface.HandFunc

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
		// 将最大的数据读到 buf 中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			return
		}
		// 调用当前连接业务(这里执行的是当前 conn 绑定的 handle 方法)
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnID, " handle is error")
			c.ExitBuffChan <- true
			continue
		}
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

func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		handleAPI:    callback_api,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
	}

	return c
}
