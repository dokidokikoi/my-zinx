package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dokidokikoi/my-zinx/ziface"
)

// 全局参数
type GlobalObj struct {
	// 当前 Zinx 全局 Server 对象
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string
	Version   string

	// 读取数据包的最大值
	MaxPacketSize uint32
	// 当前服务器主机允许的最大连接个数
	MaxConn int
	// 业务工作池的数量
	WorkerPoolSize uint32
	// 业务工作 worker 对应任务队列的最大任务存储数量
	MaxWorkerTaskLen uint32

	ConfFilePath string
}

var GlobalObject *GlobalObj

// 加载用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("加载默认配置")
		return
	}

	// 将 json 数据解析到 struct 中
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	// 初始化 GlobalObject 变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "v0.4",
		TcpPort:          7777,
		Host:             "0.0.0.0",
		MaxPacketSize:    4096,
		ConfFilePath:     "conf/zinx.json",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 从配置文件加载用户配置
	GlobalObject.Reload()
}
