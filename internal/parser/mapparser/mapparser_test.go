package mapparser

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tealeg/xlsx"
	sw "parser_server/internal/parser/sheetwrapper"
	"reflect"
	"testing"
)

func TestDefaultColumnExtractor_GetHeaders_Integration(t *testing.T) {
	// Загрузка тестового файла
	file, err := xlsx.OpenFile("test.xlsx")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	// Получение листа
	sheet := sw.DefaultSheetWrapper{
		Sheet: file.Sheets[0],
	}

	// Инициализация DefaultColumnExtractor
	dce := &DefaultColumnExtractor{
		Sheet: &sheet,
	}
	result := dce.GetHeaders()

	expected := map[string]int{
		"Some header1": 0,
		"Some header2": 1,
		"Some header3": 2,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractColumn(0) = %v, expected %v", result, expected)
	}
}

func TestDefaultColumnExtractor_ExtractColumn_Integration(t *testing.T) {
	// Загрузка тестового файла
	file, err := xlsx.OpenFile("test.xlsx")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Получение листа
	sheet := sw.DefaultSheetWrapper{
		Sheet: file.Sheets[0],
	}

	// Инициализация DefaultColumnExtractor
	dce := &DefaultColumnExtractor{
		Sheet: &sheet,
	}

	// Тестирование
	result := dce.ExtractColumn(0)
	expected := map[int]string{
		1: "Value1",
		2: "Value2",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractColumn(0) = %v, expected %v", result, expected)
	}
}

func TestTableProcessor_ExtractTable(t *testing.T) {
	// Загрузка тестового файла
	file, err := xlsx.OpenFile("test.xlsx")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	// Получение листа
	sheet := sw.DefaultSheetWrapper{
		Sheet: file.Sheets[0],
	}

	// Инициализация DefaultColumnExtractor
	dce := &DefaultColumnExtractor{
		Sheet: &sheet,
	}
	processor := NewTableProcessor(dce)

	result := processor.ExtractTable()

	expected := map[string]map[int]string{
		"Some header1": {1: "Value1", 2: "Value2"},
		"Some header2": {1: "Value21", 2: "Value22"},
		"Some header3": {1: "Value31", 2: "Value32"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractColumn(0) = %v, expected %v", result, expected)
	}
}

// TestDefaultColumnExtractor_ExtractColumn тестирует метод ExtractColumn
func TestDefaultColumnExtractor_ExtractColumn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для xlsx.Sheet
	mockSheet := sw.NewMockSheetWrapper(ctrl)

	// Настраиваем ожидания для мока
	rows := []*xlsx.Row{
		{Cells: []*xlsx.Cell{{Value: "Header1"}, {Value: "Header2"}, {Value: "Header3"}}},
		{Cells: []*xlsx.Cell{{Value: "Row1Col1"}, {Value: "Row1Col2"}, {Value: "Row1Col3"}}},
		{Cells: []*xlsx.Cell{{Value: "Row2Col1"}, {Value: "Row2Col2"}, {Value: "Row2Col3"}}},
		{Cells: []*xlsx.Cell{{Value: "Row3Col1"}, {Value: "Row3Col2"}, {Value: "Row3Col3"}}},
	}
	mockSheet.EXPECT().Rows().Return(rows).AnyTimes()

	// Создаем экземпляр DefaultColumnExtractor с моком
	extractor := &DefaultColumnExtractor{Sheet: mockSheet}

	// Тестируем извлечение колонки с индексом 1 (Header2)
	columnData := extractor.ExtractColumn(1)

	// Проверяем результат
	expected := map[int]string{
		1: "Row1Col2",
		2: "Row2Col2",
		3: "Row3Col2",
	}
	assert.Equal(t, expected, columnData)
}

// TODO: Такой же тест на моках но для GetHeaders
func TestDefaultColumnExtractor_GetHeaders(t *testing.T) {
	ctlr := gomock.NewController(t)
	defer ctlr.Finish()

	mockSheet := sw.NewMockSheetWrapper(ctlr)

	rows := []*xlsx.Row{
		{Cells: []*xlsx.Cell{{Value: "Header1"}, {Value: "Header2"}, {Value: "Header3"}}},
	}

	mockSheet.EXPECT().Rows().Return(rows).AnyTimes()

	// Создаем экземпляр DefaultColumnExtractor с моком
	extractor := &DefaultColumnExtractor{Sheet: mockSheet}

	// Тестируем извлечение колонки с индексом 1 (Header2)
	headers := extractor.GetHeaders()

	// Проверяем результат
	expected := map[string]int{
		"Header1": 0,
		"Header2": 1,
		"Header3": 2,
	}
	assert.Equal(t, expected, headers)
}
