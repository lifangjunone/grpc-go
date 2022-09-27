package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	simple "grpc-go/simple/pb"
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
}
