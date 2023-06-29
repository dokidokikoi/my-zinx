package main

import (
	"fmt"

	"github.com/dokidokikoi/my-zinx/ziface"
	"github.com/dokidokikoi/my-zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")

	_, err := request.GetConnetcion().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")

	_, err := request.GetConnetcion().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")

	_, err := request.GetConnetcion().GetTCPConnection().Write([]byte("After ping...\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {
	s := znet.NewServer("[Zinx v0.3]")

	s.AddRouter(&PingRouter{})

	s.Serve()
}
