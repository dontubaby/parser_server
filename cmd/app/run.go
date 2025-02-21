package app

import (
	"fmt"
	"log"
	"net"
	"parser_server/internal/config"
	"parser_server/server"

	"google.golang.org/grpc"
	"parser_server/server/pb"
)

func Run(cfg config.Config) error {
	//устанавливаем соединение
	lis, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()
	// Регистрируем сервис парсинга в gRPC
	pb.RegisterParseServiceServer(s, server.NewServer(cfg))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
