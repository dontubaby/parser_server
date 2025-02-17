package main

import (
	"log"
)

// структура таблицы ключ - название колонки значение -мапа ключ - номер строки значение ячейки - строка

func main() {
	// Запуск приложения
	if err := Run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}

}
