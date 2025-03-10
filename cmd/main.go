package main

import (
	"fmt"
	"log"
	"parser_server/cmd/app"
	"parser_server/internal/config"
)

func main() {
	cfg, err := config.ParseConfigFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Config data: ", cfg.FileName, cfg.ServerAddress)
	// Запуск приложения
	if err := app.Run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}

}
