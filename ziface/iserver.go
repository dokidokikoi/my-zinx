package ziface

// 服务接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务
	Stop()
	// 开启业务服务
	Serve()
}
