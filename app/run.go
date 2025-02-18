package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "parser_server/server/pb"
)

type Server struct {
	pb.UnimplementedParseServiceServer
}

func Run(cfg Config) error {
	//устанавливаем соединение
	lis, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()
	pb.RegisterParseServiceServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
