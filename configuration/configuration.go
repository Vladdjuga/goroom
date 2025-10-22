package configuration

import (
	"encoding/json"
	"os"
)

type Config struct {
	HttpPort string `json:"httpPort"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	// Попробовать загрузить из файла
	data, err := os.ReadFile(path)
	if err != nil {
		// Если файла нет (например, на Render) - использовать ENV переменные
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		config.HttpPort = ":" + port
		return &config, nil
	}

	// Если файл есть - парсим его
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Переопределить порт из ENV если задан (приоритет ENV)
	if envPort := os.Getenv("PORT"); envPort != "" {
		config.HttpPort = ":" + envPort
	}

	return &config, nil
}
