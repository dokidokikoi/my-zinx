package znet

import "github.com/dokidokikoi/my-zinx/ziface"

type Request struct {
	// 已经和客户端建立好连接
	conn ziface.IConnection
	// 客户端请求的数据
	data []byte
}

// 获取请求连接的数据
func (r *Request) GetConnetcion() ziface.IConnection {
	return r.conn
}

// 获取请求连接的数据
func (r *Request) GetData() []byte {
	return r.data
}
