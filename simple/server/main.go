package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go/middleware/server"
	simple "grpc-go/simple/pb"
	"io"
	"log"
	"net"
)

type HelloServiceServer struct {
	simple.UnimplementedHelloServiceServer
}

func (h *HelloServiceServer) Hello(ctx context.Context, req *simple.Request) (*simple.Response, error) {
	return &simple.Response{Value: fmt.Sprintf("hello %s", req.Value)}, nil
}

func (h *HelloServiceServer) Channel(stream simple.HelloService_ChannelServer) error {
	// loop receive stream
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("channel closed %s", err.Error())
				return nil
			}
			return err
		}
		log.Printf("Server receive stream: %s", req.Value)
		// loop send stream
		err = stream.Send(&simple.Response{Value: fmt.Sprintf("hello %s", req.Value)})
		if err != nil {
			if err == io.EOF {
				log.Println("channel closed")
				return nil
			}
			return err
		}
	}
}

func main() {
	// create a grpc server
	unaryInter := grpc.UnaryInterceptor(server.NewAuthUnaryInterceptor())
	streamInter := grpc.StreamInterceptor(server.NewAuthStreamInterceptor())
	grpcSvc := grpc.NewServer(
		unaryInter,
		streamInter,
	)
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
