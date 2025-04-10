package main

import (
	"fmt"
	"log"
	"parser_server/cmd/app"
	"parser_server/internal/config"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.ParseConfigFile("./config.json") //TODO: переделать на переменную окружения
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server start working: %v\n", cfg.ServerAddress)

	// Запуск приложения
	if err := app.Run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
