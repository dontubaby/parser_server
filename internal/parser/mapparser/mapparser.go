package mapparser

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	dialog "parser_server/internal/models/dialog"
	models "parser_server/internal/models/map"
	sw "parser_server/internal/parser/sheetwrapper"
	"strconv"
	"strings"
)

// ColumnExtractor Интерфейс для извлечения колонок
type ColumnExtractor interface {
	ExtractColumn(colIndex int) map[int]string //Экстрактор колонок
	GetHeaders() map[string]int                //Экстрактор заголовков
}

// DefaultColumnExtractor Объект реализующий интерфейс извлечения колонок таблицы
type DefaultColumnExtractor struct {
	Sheet sw.SheetWrapper
}

func NewDefaultColumnExtractor(sheet *xlsx.Sheet) *DefaultColumnExtractor {
	return &DefaultColumnExtractor{
		Sheet: sw.NewDefaultSheetWrapper(sheet),
	}
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

func (dce *DefaultColumnExtractor) GetHeaders() map[string]int {
	rows := dce.Sheet.Rows()
	if len(rows) == 0 {
		return nil
	}
	headers := make(map[string]int)
	for i, cell := range rows[0].Cells {
		headers[cell.String()] = i
	}
	return headers
}

func (dce *DefaultColumnExtractor) ParseTasks(filename string) ([]dialog.Task, error) {
	var tasks []dialog.Task
	log.Printf("File %s opened", filename)
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	sheet := file.Sheet["Задания"]
	if sheet == nil {
		return nil, fmt.Errorf("sheet 'Задания' not found")
	}

	headers := dce.GetHeaders()
	for _, row := range sheet.Rows[1:] { // Пропускаем заголовок
		task, err := dce.ParseRowToTask(row.Cells, headers)
		//fmt.Printf("TASK: %v\n", task)
		if err != nil {
			continue // Пропускаем проблемные строки
		}
		tasks = append(tasks, *task)

	}

	return tasks, nil
}

func (dce *DefaultColumnExtractor) ParseRowToTask(cells []*xlsx.Cell, headers map[string]int) (*dialog.Task, error) {
	task := dialog.Task{}
	//C1
	profession := cells[headers["Профессия"]].String()
	if profession == "" {
		return nil, fmt.Errorf("missing Profession")
	}
	//C2
	taskIDStr := cells[headers["ID задания"]].String()
	if taskIDStr == "" {
		return nil, fmt.Errorf("missing TaskID")
	}
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TaskID: %v", err)
	}
	//C3
	taskGroupIDStr := cells[headers["ID группы заданий"]].String()
	if taskGroupIDStr == "" {
		return nil, fmt.Errorf("missing TaskGroupID")
	}
	taskGroupID, err := strconv.ParseUint(taskGroupIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TaskGroupID: %v", err)
	}
	//C4
	groupWeightStr := cells[headers["Вес группы в ленте профессии (очередь)"]].String()
	if groupWeightStr == "" {
		return nil, fmt.Errorf("missing GroupWeight")
	}
	groupWeight, err := strconv.ParseUint(groupWeightStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid GroupWeight: %v", err)
	}
	//C5
	positionStr := cells[headers["Позиция в группе"]].String()
	if positionStr == "" {
		return nil, fmt.Errorf("missing Position")
	}
	position, err := strconv.ParseUint(positionStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid Position: %v", err)
	}
	//C6
	xpProfessionStr := cells[headers["XP профессии за задание"]].String()
	if xpProfessionStr == "" {
		return nil, fmt.Errorf("missing ProfessionXP")
	}
	xpProfession, err := strconv.ParseUint(xpProfessionStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid ProfessionXP: %v", err)
	}
	//C7
	tmpXPStr := cells[headers["Временное XP профессии за задание"]].String()
	if tmpXPStr == "" {
		tmpXPStr = "0" // Устанавливаем значение по умолчанию
	}
	tmpXP, err := strconv.ParseUint(tmpXPStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TmpXP: %v", err)
	}
	//C8
	playerXPStr := cells[headers["XP игрока за задание"]].String()
	if playerXPStr == "" {
		return nil, fmt.Errorf("missing PlayerXp")
	}
	playerXP, err := strconv.ParseUint(playerXPStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid PlayerXp: %v", err)
	}
	//C9
	activationXPStr := cells[headers["XP профессии для активации группы (не временные)"]].String()
	if activationXPStr == "" {
		activationXPStr = "0" // Устанавливаем значение по умолчанию
	}
	activationXP, err := strconv.ParseUint(activationXPStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid ActivationXP: %v", err)
	}
	//C10
	dialogID := cells[headers["ID диалогов"]].String()
	if dialogID == "" {
		dialogID = "" // Может быть пустым
	}
	//C11
	activationAct := cells[headers["активация группы активирует Акт"]].String()
	if activationAct == "" {
		activationAct = "" // Может быть пустым
	}
	//C12
	nextActXPStr := cells[headers["XP для перехода в следующий Акт"]].String()
	if nextActXPStr == "" {
		nextActXPStr = "0" // Устанавливаем значение по умолчанию
	}
	nextActXP, err := strconv.ParseUint(nextActXPStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid NextActXP: %v", err)
	}
	//C13
	nextChapterXPStr := cells[headers["XP для перехода в следующую Главу"]].String()
	if nextChapterXPStr == "" {
		nextChapterXPStr = "0" // Устанавливаем значение по умолчанию
	}
	nextChapterXP, err := strconv.ParseUint(nextChapterXPStr, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid NextChapterXP: %v", err)
	}
	//TODO: Проверить необходимость приведения типов
	task = dialog.Task{
		Profession:    profession,
		TaskID:        taskID,
		TaskGroupID:   taskGroupID,
		GroupWeight:   uint32(groupWeight),
		Position:      uint32(position),
		ProfessionXP:  uint32(xpProfession),
		TmpXP:         uint32(tmpXP),
		PlayerXp:      uint32(playerXP),
		ActivationXP:  uint32(activationXP),
		DialogID:      dialogID,
		ActivationAct: activationAct,
		NextActXP:     uint32(nextActXP),
		NextChapterXP: uint32(nextChapterXP),
	}
	fmt.Printf("TASK: %v\n", task)
	return &task, nil
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

	for headerName, colIndex := range headers {
		columnData := tp.extractor.ExtractColumn(colIndex)
		if len(columnData) > 0 {
			table[headerName] = columnData // Используем строковый ключ
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
