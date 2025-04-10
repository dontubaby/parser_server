package parser

import (
	"fmt"
	"github.com/tealeg/xlsx"
	models "parser_server/internal/models/map"
	"strconv"
	"strings"
)

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
	//i - номер строки, dce.Sheet.Rows[i].Cells[colIndex].String() - значение ячейки искомой колонки в этой строке
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

// Функция для парсинга координат
func parseCoordinate(coord string) [2]int {
	parts := strings.Split(coord, ".")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return [2]int{x, y}
}

func parseExcelToMap(filePath string) (models.AutobattlerMap, error) {
	var parsedMap models.AutobattlerMap

	// Открытие файла Excel
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return parsedMap, fmt.Errorf("error opening file: %w", err)
	}

	// Получение листа "Автобатлер"
	sheet := file.Sheet["Автобатлер"]
	if sheet == nil {
		return parsedMap, fmt.Errorf("sheet 'Автобатлер' not found")
	}

	// Мапы для быстрого доступа к секторам и зонам
	sectorMap := make(map[string]*models.Sector)
	zoneMap := make(map[string]*models.Zone)

	// Проход по строкам таблицы
	for rowIndex, row := range sheet.Rows {
		if rowIndex == 0 { // Пропускаем заголовок
			continue
		}

		// Чтение ячеек строки
		cells := row.Cells
		if len(cells) < 6 || cells[0].String() == "" {
			continue // Пропускаем пустые строки
		}

		sectorID := cells[1].String()
		zoneID := cells[3].String()
		hexCoordStr := cells[2].String()
		hexType := cells[5].String()
		hexModifiers := cells[4].String()

		// Проверка наличия обязательных данных
		if sectorID == "" || zoneID == "" || hexCoordStr == "" || hexType == "" {
			continue
		}

		// Создание гекса
		hex := models.Hex{
			Coordinate: parseCoordinate(hexCoordStr),
			Htype:      hexType,
		}

		// Быстрый доступ к сектору через мапу
		sector, exists := sectorMap[sectorID]
		if !exists {
			// Если сектора нет, создаём новый
			newSector := models.Sector{Id: sectorID, Zones: []models.Zone{}}

			parsedMap.Sectors = append(parsedMap.Sectors, newSector)
			sectorMap[sectorID] = &parsedMap.Sectors[len(parsedMap.Sectors)-1]
			sector = sectorMap[sectorID]
		}

		// Быстрый доступ к зоне через мапу
		zone, zoneExists := zoneMap[zoneID]
		if !zoneExists {
			// Если зоны нет, создаём новую
			newZone := models.Zone{
				Id:        zoneID,
				Modifiers: hexModifiers,
				Hexs:      []models.Hex{},
			}
			sector.Zones = append(sector.Zones, newZone)
			zoneMap[zoneID] = &sector.Zones[len(sector.Zones)-1]
			zone = zoneMap[zoneID]
		}

		// Добавление гекса в зону
		zone.Hexs = append(zone.Hexs, hex)
	}

	return parsedMap, nil
}
