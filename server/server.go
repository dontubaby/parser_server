package server

import (
	"context"
	"github.com/tealeg/xlsx"
	"log"

	"parser_server/internal/config"
	"parser_server/internal/parser"
	pb "parser_server/server/pb"
)

const (
	interval = 10
)

type Server struct {
	pb.UnimplementedParseServiceServer
	parser parser.TableProcessor
}

func (s *Server) GetTable(ctx context.Context) (*pb.Table, error) {
	cfg, err := config.ParseConfigFile("config.json")
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}

	xlFile, err := xlsx.OpenFile(cfg.FileName)
	if err != nil {
		log.Fatalf("open file error: %v", err)
	}
	sheet := xlFile.Sheets[0]

	// Создаем экстрактор и процессор
	extractor := &parser.DefaultColumnExtractor{Sheet: sheet}
	processor := parser.NewTableProcessor(extractor)

	tableData := processor.ExtractTable()

	// Преобразуем таблицу в protobuf-формат
	pbTable := &pb.Table{
		Table: make(map[string]*pb.TableRows),
	}
	for header, column := range tableData {
		pbRow := &pb.TableRows{
			Rows: make(map[int32]string),
		}
		for rowIndex, value := range column {
			pbRow.Rows[int32(rowIndex)] = value
		}
		pbTable.Table[header] = pbRow
	}

	return pbTable, nil
}
