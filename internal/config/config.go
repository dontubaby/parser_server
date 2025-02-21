package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address"`
	FileName      string `json:"file_name"`
}

func ParseConfigFile(filename string) (Config, error) {
	var data []byte
	configFile, err := os.Open(filename)
	if err != nil {
		log.Printf("Open file error - %v", err)
		return Config{}, err
	}
	//Запускаем закрытие файла из замыкания, чтобы была возможность обработать ошибку закрытия файла
	defer func() {
		err = configFile.Close()
		if err != nil {
			log.Printf("Close file error - %v", err)
		}
	}() // _ = configFile.Close()

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	jsonErr := json.Unmarshal(data, &config)
	if jsonErr != nil {
		log.Printf("Unmarhaling error - %v", err)
		return Config{}, err
	}
	return config, nil
}
