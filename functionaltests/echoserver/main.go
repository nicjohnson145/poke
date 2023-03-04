package main

import (
	"context"
	"fmt"
	"net"
	"os"

	pb "github.com/nicjohnson145/poke/functionaltests/echoserver/protobuf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (e *EchoSvr) Err(_ context.Context, req *pb.ErrMsg) (*emptypb.Empty, error) {
	code := codes.Unknown
	if req.Code != 0 {
		code = codes.Code(req.Code)
	}
	return nil, status.Error(code, "you did bad")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	check(err)

	logger := InitLogger()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(DefaultMethodLoggingInterceptor(logger)))
	reflection.Register(server)
	pb.RegisterEchoServiceServer(server, &EchoSvr{})

	fmt.Println("starting server")
	if err := server.Serve(lis); err != nil {
		die(err.Error())
	}
}

func DefaultMethodLoggingInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Info().Str("path", info.FullMethod).Msg("request recieved")
		return handler(ctx, req)
	}
}

func InitLogger() zerolog.Logger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	return log.With().Logger()
}
