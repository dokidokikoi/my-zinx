package ziface

type IMsgHandler interface {
	// 马上以非阻塞的方式处理消息
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
}
