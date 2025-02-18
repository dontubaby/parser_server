package parser

import "github.com/tealeg/xlsx"

// ColumnExtractor Интерфейс для извлечения колонок
type ColumnExtractor interface {
	ExtractColumn(colIndex int) map[int]string //Экстрактор колонок
	GetHeaders() []string                      //Экстрактор заголовков
}

// DefaultColumnExtractor Объект реализующий интерфейс извлечения колонок таблицы
type DefaultColumnExtractor struct {
	Sheet *xlsx.Sheet
}

func (dce *DefaultColumnExtractor) ExtractColumn(colIndex int) map[int]string {
	columnValues := make(map[int]string)
	for rowIndex, row := range dce.Sheet.Rows {
		if len(row.Cells) > colIndex {
			columnValues[rowIndex] = row.Cells[colIndex].String()
		}
	}
	return columnValues
}

func (dce *DefaultColumnExtractor) GetHeaders() []string {
	if len(dce.Sheet.Rows) == 0 {
		return nil
	}
	headers := make([]string, len(dce.Sheet.Rows[0].Cells))
	for i, cell := range dce.Sheet.Rows[0].Cells {
		headers[i] = cell.String()
	}
	return headers
}

type TableProcessor struct {
	extractor ColumnExtractor
}

func NewTableProcessor(extractor ColumnExtractor) *TableProcessor {
	return &TableProcessor{extractor: extractor}
}

func (tp *TableProcessor) ExtractTable() map[string]map[int]string {
	table := make(map[string]map[int]string)
	headers := tp.extractor.GetHeaders()

	for i := 0; i < len(headers); i++ {
		columnData := tp.extractor.ExtractColumn(i)
		if len(columnData) > 0 {
			header := headers[i]
			table[header] = columnData
		}
	}
	return table
}
