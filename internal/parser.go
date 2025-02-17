package internal

import "github.com/tealeg/xlsx"

type Table struct {
	data map[string]map[int]string
}

func Extractor(colIndex int, sheet *xlsx.Sheet) map[int]string {
	columnValues := make(map[int]string)
	for rowIndex, row := range sheet.Rows {
		if len(row.Cells) > colIndex {
			columnValues[rowIndex] = row.Cells[colIndex].String()
		}
	}

	return columnValues
}

func TableExtractor(sheet *xlsx.Sheet) map[string]map[int]string {
	table := make(map[string]map[int]string)
	headers := sheet.Rows[0].Cells
	for i := 0; i <= len(headers); i++ {
		columnData := Extractor(i, sheet)
		if len(columnData) > 0 {
			header := headers[i].String()
			table[header] = columnData
		}
	}

	return table
}
