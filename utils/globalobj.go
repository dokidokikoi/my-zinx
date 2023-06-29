package utils

import (
	"encoding/json"
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
}

var GlobalObject *GlobalObj

// 加载用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
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
		Name:          "ZinxServerApp",
		Version:       "v0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxPacketSize: 4096,
	}

	// 从配置文件加载用户配置
	GlobalObject.Reload()
}
