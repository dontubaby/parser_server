package server

import (
	"context"
	"fmt"
	"github.com/tealeg/xlsx"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	dialog "parser_server/internal/models/dialog"
	"parser_server/internal/parser/mapparser"

	"parser_server/internal/config"
	"parser_server/server/pb"
)

type Server struct {
	pb.UnimplementedParseAutoBattlerServiceServer
	pb.UnimplementedParseDialogServiceServer
	parser *mapparser.TableProcessor
	cfg    config.Config
}

func NewServer(cfg config.Config, extractor mapparser.ColumnExtractor) *Server {
	return &Server{
		parser: mapparser.NewTableProcessor(extractor),
		cfg:    cfg,
	}
}

func (s *Server) GetAutobattlerTable(ctx context.Context, req *emptypb.Empty) (*pb.AutobattlerTable, error) {
	xlFile, err := xlsx.OpenFile("./" + s.cfg.AutoBattlerConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	log.Printf("File %s opened", s.cfg.AutoBattlerConfigFileName)
	sheet := xlFile.Sheets[0]
	// Создаем экстрактор и процессор
	extractor := mapparser.NewDefaultColumnExtractor(sheet)
	processor := mapparser.NewTableProcessor(extractor)
	tableData := processor.ExtractTable()
	// Преобразуем таблицу в protobuf-формат
	pbTable := &pb.AutobattlerTable{
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

// TODO:написать тесты
func ConvertToPbTaskList(tasks []dialog.Task) *pb.TaskList {
	pbTasks := &pb.TaskList{
		Table: make([]*pb.Task, len(tasks)),
	}
	for i, task := range tasks {
		pbTasks.Table[i] = &pb.Task{
			Profession:    task.Profession,
			TaskID:        task.TaskID,
			TaskGroupID:   task.TaskGroupID,
			GroupWeight:   task.GroupWeight,
			Position:      task.Position,
			ProfessionXP:  task.ProfessionXP,
			TmpXP:         task.TmpXP,
			PlayerXp:      task.PlayerXp,
			ActivationXP:  task.ActivationXP,
			DialogID:      task.DialogID,
			ActivationAct: task.ActivationAct,
			NextActXP:     task.NextActXP,
			NextChapterXP: task.NextChapterXP,
		}
	}

	return pbTasks
}

func (s *Server) GetDialogTable(ctx context.Context, req *emptypb.Empty) (*pb.TaskList, error) {
	xlFile, err := xlsx.OpenFile("./" + s.cfg.DialogConfigFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	log.Printf("File %s opened", s.cfg.DialogConfigFileName)
	sheet := xlFile.Sheets[2] //2- это нимер листа 'Задания' TODO: убери HARDCODE
	fmt.Printf("Sheet name: %s", sheet.Name)
	// Создаем экстрактор
	extractor := mapparser.NewDefaultColumnExtractor(sheet)
	//Извлекаем задачи из конфигурационного файла системы диалогов
	tasks, err := extractor.ParseTasks("./" + s.cfg.DialogConfigFileName)
	fmt.Println(tasks)
	if err != nil {
		return nil, fmt.Errorf("error parsing tasks: %w", err)
	}

	// Преобразуем массив dialog.Task в pb.TaskList
	pbTasks := ConvertToPbTaskList(tasks)

	return pbTasks, nil
}
