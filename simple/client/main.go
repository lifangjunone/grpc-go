package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
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
	rsp, err := client.Hello(context.Background(), &simple.Request{Value: "ldd"})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Value)
	// stream example
	stream, err := client.Channel(context.Background())
	if err != nil {
		panic(err)
	}
	// send stream
	go func() {
		for {
			if err := stream.Send(&simple.Request{Value: "ldd"}); err != nil {
				if err == io.EOF {
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
