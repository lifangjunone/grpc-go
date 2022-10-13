package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	client2 "grpc-go/middleware/client"
	simple "grpc-go/simple/pb"
	"io"
	"log"
	"time"
)

func main() {
	// 为客户端添加认证信息
	credentialInfo := grpc.WithPerRPCCredentials(client2.NewAuthentication("admin", "123456"))
	// 建立网络链接
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:1234", grpc.WithInsecure(), credentialInfo)
	if err != nil {
		panic(err)
	}
	client := simple.NewHelloServiceClient(conn)
	// md := server.NewClientAuth("admin", "123456")
	// ctx := metadata.NewOutgoingContext(context.Background(), md)
	rsp, err := client.Hello(context.Background(), &simple.Request{Value: "ldd"})
	if err != nil {
		fmt.Printf("call is failed %#v", err.Error())
		return
	}
	fmt.Println(rsp.Value)
	// stream example
	stream, err := client.Channel(context.Background())
	if err != nil {
		fmt.Printf("call is failed %#v", err.Error())
	}
	// send stream
	go func() {
		for {
			if err := stream.Send(&simple.Request{Value: "ldd"}); err != nil {
				if err == io.EOF {
					fmt.Printf("send info to server")
					log.Println("server closed")
				}
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()
	// receive steam
	for {
		rsp, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println("server closed")
			}
			return
		}
		log.Printf("Client receive stream %s", rsp.Value)
		time.Sleep(2 * time.Second)
	}
}
