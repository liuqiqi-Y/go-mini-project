package main

import (
	"context"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type HelloServiceImpl struct {
	UnimplementedHelloServiceServer
}

func (p *HelloServiceImpl) Hello(ctx context.Context, args *String) (*String, error) {
	reply := &String{Value: "hello " + args.GetValue()}
	return reply, nil
}
func (p *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		/*
			服务端在循环中接收客户端发来的数据，
			如果遇到io.EOF表示客户端流被关闭，
			如果函数退出表示服务端流关闭。
			生成返回的数据通过流发送给客户端，
			双向流数据的发送和接收都是完全独立的行为。
			需要注意的是，发送和接收的操作并不需要一一对应，
			用户可以根据真实场景进行组织代码。
		*/
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		reply := &String{Value: "hello " + args.GetValue()}
		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}
func main() {
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))
	lis, err := net.Listen("tcp", ":8899")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
