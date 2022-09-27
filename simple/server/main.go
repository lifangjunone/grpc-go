package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	simple "grpc-go/simple/pb"
	"log"
	"net"
)

type HelloServiceServer struct {
	simple.UnimplementedHelloServiceServer
}

func (h *HelloServiceServer) Hello(ctx context.Context, req *simple.Request) (*simple.Response, error) {
	return &simple.Response{Value: fmt.Sprintf("hello %s", req.Value)}, nil
}

func main() {
	// create a grpc server
	grpcSvc := grpc.NewServer()
	// register HelloServer to grpc service
	simple.RegisterHelloServiceServer(grpcSvc, new(HelloServiceServer))
	// listen server
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	log.Printf("Server start ...")
	log.Printf("Listen: http://127.0.0.1:1234")
	// start grpc server
	if err = grpcSvc.Serve(listener); err != nil {
		panic(err)
	}

}
