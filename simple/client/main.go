package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc-go/middleware/server"
	simple "grpc-go/simple/pb"
	"io"
	"log"
	"time"
)

func main() {
	// 建立网络链接
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:1234", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := simple.NewHelloServiceClient(conn)
	md := server.NewClientAuth("admin", "123456")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	rsp, err := client.Hello(ctx, &simple.Request{Value: "ldd"})
	if err != nil {
		fmt.Printf("call is failed %#v", err.Error())
		return
	}
	fmt.Println(rsp.Value)
	// stream example
	stream, err := client.Channel(ctx)
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
