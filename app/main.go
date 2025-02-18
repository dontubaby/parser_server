package main

import (
	"log"
)

var cfg = Config{
	ServerAddress: "localhost:8080",
	FileName:      "cyconfig.xlsx",
}

func main() {
	// Запуск приложения
	if err := Run(cfg); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}

}
