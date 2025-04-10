package parser

import "github.com/tealeg/xlsx"

// SheetWrapper это интерфейс-обертка вокруг xlsx.Sheet
type SheetWrapper interface {
	Rows() []*xlsx.Row
}

// DefaultSheetWrapper это реализация интерфейса SheetWrapper
type DefaultSheetWrapper struct {
	Sheet *xlsx.Sheet
}

// Rows возвращает строки из xlsx.Sheet
func (dsw *DefaultSheetWrapper) Rows() []*xlsx.Row {
	return dsw.Sheet.Rows
}
