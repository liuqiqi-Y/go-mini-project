package main

import (
	"fmt"
	"net/rpc"
)

const HelloServiceName = "path/to/pkg.HelloService"

type HelloServiceInterface = interface {
	Hello(request string, reply *string) error
}
type HelloServiceClient struct {
	*rpc.Client
}

func (h *HelloServiceClient) Hello(request string, reply *string) error {
	return h.Call(HelloServiceName+".Hello", request, reply)
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(net, address string) (*HelloServiceClient, error) {
	conn, err := rpc.Dial(net, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{conn}, nil
}

func main() {
	conn, err := DialHelloService("tcp", "localhost:8899")
	if err != nil {
		panic("建立连接失败: " + err.Error())
	}
	var reply string
	err = conn.Hello("WANG", &reply)
	if err != nil {
		fmt.Printf("请求远程服务失败: %s\n", err.Error())
		return
	}
	fmt.Println(reply)
}
