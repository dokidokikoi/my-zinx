package znet

import (
	"fmt"
	"github.com/dokidokikoi/my-zinx/ziface"
	"strconv"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgID=%d is not FOUND!\n", request.GetMsgID())
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.判断当前 msg 绑定的 API 处理方式是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2.添加 msg 与 api 的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
