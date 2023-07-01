package znet

import (
	"errors"
	"fmt"
	"github.com/dokidokikoi/my-zinx/utils"
	"io"
	"net"
	"sync"

	"github.com/dokidokikoi/my-zinx/ziface"
)

type Connection struct {
	// 当前 Conn 属于哪个 Server
	TcpServer ziface.IServer
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

	// 无缓冲管道，用于读、写两个 Goroutine 之间的数据通信
	msgChan chan []byte
	// 有缓冲管道，用于读、写两个 Goroutine 之间的数据通信
	msgBuffChan chan []byte
	// 连接属性
	property map[string]interface{}
	// 保护连接属性修改的锁
	propertyLock sync.RWMutex
}

// 停止连接，结束当前连接状态 M
func (c *Connection) Stop() {
	// 1. 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 如果用户注册了该连接的关闭回调业务，
	// 则在此时应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭 Socket 连接
	c.Conn.Close()

	// 通知从缓冲队列读数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	// 将连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭该连接的全部的管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
	close(c.msgChan)
}

// 处理 conn 读数据的 Gorutine
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.Conn.RemoteAddr().String(), "[Conn Reader exit!]")
	defer c.Stop()

	for {
		// 创建拆包对象
		dp := NewDataPack()

		// 读取客户端的 Msg Head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			return
		}
		// 拆包
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			return
		}

		// 根据 dataLen 读取数据
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				return
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经启动 worker 工作池，将消息交给 worker
			c.MsgHandler.SendMsg2TaskQueue(&req)
		} else {
			// 从绑定好的消息和对应的处理方法中执行对应的 Handle 方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	// 开启处理该连接读取客户端数据的 Goroutine
	go c.StartReader()
	// 开启用于写回客户端的数据流程的 Goroutine
	go c.StartWriter()

	// 执行用户传入的构造方法
	c.TcpServer.CallOnConnStart(c)

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
	// 将得到的数据发送到 chan，供 writer 读取
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
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
	// 将得到的数据发送到 chan，供 writer 读取
	c.msgBuffChan <- msg

	return nil
}

// 读写分离，职责单一，在优化读或写逻辑时互不干扰
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[Conn Writer exit!]")
	defer c.Stop()

	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Printf("Send Data error: %v, Conn Writer exit", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			// 针对有缓冲的 chan 需要进行数据处理
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Printf("Send Data error: %v, Conn Writer exit", err)
					return
				}
			}
			break
			fmt.Println("msgBuffChan is Closed")
		case <-c.ExitBuffChan:
			// conn 已经关闭
			return
		}
	}
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	val, ok := c.property[key]
	if !ok {
		return nil, errors.New("no property found")
	}
	return val, nil
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   msgHandler,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}

	// 将新创建的 Conn 添加到连接管理中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}
