package znet

import (
	"fmt"
	"net"

	"github.com/dokidokikoi/my-zinx/ziface"
)

// IServer 的接口实现
type Server struct {
	// 服务器名
	Name string
	// tcp4 or other
	IPVersion string
	// 服务绑定的 IP 地址
	IP string
	// 服务绑定的端口
	Port int
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	// 开启一个 go 去做服务器的 listener 业务
	go func() {
		// 1.获取一个 TCP 的 Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		// 监听成功
		fmt.Println("start Zinx server", s.Name, " suc, now listening...")

		// 3.启动 server 网络连接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				return
			}

			// 3.2 TODO: Server.Start() 设置服务器最大连接控制，
			// 如果超过最大连，则关闭新的连接

			// 3.3 TODO: Server.Start() 处理该新连接1请求的业务方法，
			// 此时 handler 和 conn 应该是绑定的

			// 这类暂时做一个最大 512 字节的回显数据
			go func() {
				// 不断循环，从客户端获取数据
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					// 回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err ", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server, name", s.Name)

	// TODO: Server.Stop() 将需要清理的连接信息或者其他信息一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()

	// TODO: Server.Serve() 如果在启动服务的时候还要处理其他事情，
	// 则可以在这里添加

	// 阻塞，否则主 Go 退出，listener 的 go 将退出
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}

	return s
}
