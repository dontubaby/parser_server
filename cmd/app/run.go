package app

import (
	"fmt"
	"log"
	"net"
	"parser_server/internal/config"
	"parser_server/internal/parser/mapparser"
	"parser_server/server"

	"google.golang.org/grpc"
	"parser_server/server/pb"
)

func Run(cfg config.Config) error {
	//устанавливаем соединение
	listener, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Создаем gRPC сервер
	s := grpc.NewServer()
	// Регистрируем сервис парсинга конфигурационного файла автобатлера в gRPC
	pb.RegisterParseAutoBattlerServiceServer(s, server.NewServer(cfg, &mapparser.DefaultColumnExtractor{}))

	//Регистрируем сервис парсинга конфигурационного файла диалоговой системы в gRPC
	pb.RegisterParseDialogServiceServer(s, server.NewServer(cfg, &mapparser.DefaultColumnExtractor{}))

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
