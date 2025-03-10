package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"parser_server/internal/config"
	"parser_server/internal/parser"
	"parser_server/server/pb"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	type args struct {
		parser parser.ColumnExtractor
		cfg    config.Config
	}
	extractor := &parser.DefaultColumnExtractor{}

	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "TestNewServer",
			args: args{
				parser: extractor,
				cfg: config.Config{
					ServerAddress: "localhost:50051",
					FileName:      "test.xlsx",
				},
			},
			want: &Server{
				parser: parser.NewTableProcessor(extractor),
				cfg: config.Config{
					ServerAddress: "localhost:50051",
					FileName:      "test.xlsx",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.cfg, tt.args.parser); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTable(t *testing.T) {
	cfg, err := config.ParseConfigFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем экземпляр Server с конфигурацией
	server := &Server{
		cfg: cfg,
	}

	// Вызываем тестируемую функцию
	ctx := context.Background()
	req := &emptypb.Empty{}
	result, err := server.GetTable(ctx, req)

	// Проверяем отсутствие ошибок
	require.NoError(t, err, "Unexpected error during GetTable execution")

	// Проверяем результат
	expectedTable := &pb.Table{
		Table: map[string]*pb.TableRows{
			"Header1": {
				Rows: map[int32]string{
					1: "Value1",
					2: "Value3",
				},
			},
			"Header2": {
				Rows: map[int32]string{
					1: "Value2",
					2: "Value4",
				},
			},
		},
	}
	assert.Equal(t, expectedTable, result, "Result does not match expected table")
}
