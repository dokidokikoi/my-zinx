package main

import (
	"fmt"
	"github.com/dokidokikoi/my-zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")

	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit")
		return
	}

	for {
		// 发封包消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMessage(0, []byte("Zinx v0.6 Client Test Message0")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		msg, _ = dp.Pack(znet.NewMessage(1, []byte("Zinx v0.6 Client Test Message1")))
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head error ")
			return
		}

		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack head error ")
			return
		}
		data := make([]byte, msgHead.GetDataLen())
		if msgHead.GetDataLen() > 0 {
			_, err = io.ReadFull(conn, data)
			if err != nil {
				fmt.Println("unpack head error ")
				return
			}
		}
		msgHead.SetData(data)

		fmt.Printf("server call back Msg ID: %d, data = %s\n", msgHead.GetMsgID(), msgHead.GetData())

		time.Sleep(time.Second)
	}
}
