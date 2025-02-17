package server

import (
	"context"
	"time"

	//"github.com/tealeg/xlsx"
	pb "parser_server/server/pb"
)

const (
	interval = 10
)

type Server struct {
	pb.UnimplementedParseServiceServer
}

// структуры TableRequest и TableResponse необходимо описать в протофайле
func (s *Server) GetTable(ctx context.Context) (*pb.Table, error) {
	var table map[string]map[int]string
	table = TableExtractor(sheet)
}
