package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const HelloServiceName = "path/to/pkg.HelloService"

// 服务接口
type HelloServiceInterface = interface {
	Hello(request String, reply *String) error
}

// 服务实体
type HelloService struct{}

func (h *HelloService) Hello(request String, reply *String) error {
	reply.Value = "Hello " + request.GetValue() + "!"
	return nil
}

// 验证服务实体实现了服务接口
var _ HelloServiceInterface = (*HelloService)(nil)

// 服务注册
func RegisterHelloService(hello HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, hello)
}
func main() {
	if err := RegisterHelloService(&HelloService{}); err != nil {
		panic("注册服务失败: " + err.Error())
	}
	listener, err := net.Listen("tcp", ":8899")
	if err != nil {
		panic("开启监听端口失败: " + err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("接受客户端连接失败: " + err.Error())
		}
		//go rpc.ServeConn(conn)
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
