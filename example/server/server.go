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

func main() {
	s := znet.NewServer()

	s.AddRouter(&PingRouter{})

	s.Serve()
}
