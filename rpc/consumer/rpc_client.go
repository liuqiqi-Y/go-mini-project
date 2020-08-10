package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const HelloServiceName = "path/to/pkg.HelloService"

type HelloServiceInterface = interface {
	Hello(request *String, reply *String) error
}
type HelloServiceClient struct {
	*rpc.Client
}

func (h *HelloServiceClient) Hello(request *String, reply *String) error {
	return h.Call(HelloServiceName+".Hello", request, reply)
}

var _ HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(tcp, address string) (*HelloServiceClient, error) {
	//conn, err := rpc.Dial(net, address)
	conn, err := net.Dial(tcp, address)
	if err != nil {
		return nil, err
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{client}, nil
}

func main() {
	conn, err := DialHelloService("tcp", "localhost:8899")
	if err != nil {
		panic("建立连接失败: " + err.Error())
	}
	//var reply string
	var reply String
	request := String{Value: "WANG"}
	err = conn.Hello(&request, &reply)
	if err != nil {
		fmt.Printf("请求远程服务失败: %s\n", err.Error())
		return
	}
	fmt.Println(reply.Value)
}
