package main

import (
	"fmt"

	"github.com/dokidokikoi/my-zinx/ziface"
	"github.com/dokidokikoi/my-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

//func (pr *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PreHandle")
//
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	fmt.Printf("==> Recv Msg ID=%d, Data=%s\n", request.GetMsgID(), request.GetData())

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//func (pr *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle")
//
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
//	if err != nil {
//		fmt.Println("call back ping ping ping error")
//	}
//}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (hzr *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	fmt.Printf("==> Recv Msg ID=%d, Data=%s\n", request.GetMsgID(), request.GetData())

	err := request.GetConnection().SendMsg(1, []byte("Hello Zinx Router v0.6\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called...")
	fmt.Println("Set conn name, Home done!")
	conn.SetProperty("name", "doki")
	conn.SetProperty("home", "github.com/dokidokikoi/my-zinx")
	err := conn.SendMsg(2, []byte("DoConnection Begin..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("DoConnectionLost is Called...")

	name, _ := conn.GetProperty("name")
	home, _ := conn.GetProperty("home")
	fmt.Printf("name = %v, home = %v\n", name, home)
}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	s.Serve()
}
