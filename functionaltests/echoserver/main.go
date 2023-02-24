package main

import (
	"context"
	"fmt"
	"net"
	"os"

	pb "github.com/nicjohnson145/poke/functionaltests/echoserver/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func check(err error) {
	if err != nil {
		die(err.Error())
	}
}

func die(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type EchoSvr struct {
	pb.UnimplementedEchoServiceServer
}

func (e *EchoSvr) Echo(_ context.Context, req *pb.EchoMsg) (*pb.EchoMsg, error) {
	msg := "hello world"
	if req.Message != "" {
		msg = req.Message
	}
	return &pb.EchoMsg{Message: msg}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	check(err)

	server := grpc.NewServer()
	reflection.Register(server)
	pb.RegisterEchoServiceServer(server, &EchoSvr{})

	fmt.Println("starting server")
	if err := server.Serve(lis); err != nil {
		die(err.Error())
	}
}
