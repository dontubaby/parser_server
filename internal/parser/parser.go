package parser

// ColumnExtractor Интерфейс для извлечения колонок
type ColumnExtractor interface {
	ExtractColumn(colIndex int) map[int]string //Экстрактор колонок
	GetHeaders() []string                      //Экстрактор заголовков
}

// DefaultColumnExtractor Объект реализующий интерфейс извлечения колонок таблицы
type DefaultColumnExtractor struct {
	Sheet SheetWrapper
}

// ExtractColumn получает на вход индекс колонки и возвращает map, где ключ - номер строки, значение - значение ячейки.
func (dce *DefaultColumnExtractor) ExtractColumn(colIndex int) map[int]string {
	//Создается map для сохранения значений колонки
	columnValues := make(map[int]string)
	rows := dce.Sheet.Rows()
	//Цикл обходит колонку, пропуская 0-ю строку (с хедерами) и записывает значение ячейки в map
	//i - номер строки, dce.Sheet.Rows[i].Cells[colIndex].String() - значение ячейки искомой колонки в этой строки
	for i := 1; i < len(rows); i++ {
		if len(rows[i].Cells) > colIndex {
			columnValues[i] = rows[i].Cells[colIndex].String()
		}
	}
	return columnValues
}

func (dce *DefaultColumnExtractor) GetHeaders() []string {
	rows := dce.Sheet.Rows()
	if len(rows) == 0 {
		return nil
	}
	headers := make([]string, len(rows[0].Cells))
	for i, cell := range rows[0].Cells {
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
