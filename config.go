package main

import (
	"encoding/json"
	"os"
)

// config.json dosyasındaki değerleri tutacak yapı
type Config struct {
	BaseURL    string `json:"baseUrl"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	PageID     string `json:"page_id"`
	ReportsDir string `json:"reports_dir"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
