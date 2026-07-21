package config

import (
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port          string
	CheckInterval time.Duration
	FilePath      string
	Sites         []string
}

// LoadConfig читает переменные окружения напрямую из системы
func LoadConfig() *Config {
	// Читаем порт (если пусто, ставим 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Читаем интервал (если пусто, ставим 10 секунд)
	intervalStr := os.Getenv("CHECK_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Println("[WARN] Неверный CHECK_INTERVAL, используем 10s по умолчанию")
		interval = 10 * time.Second
	}

	// Читаем путь к файлу
	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		filePath = "metrics.json"
	}

	// Читаем сайты
	sitesStr := os.Getenv("SITES")
	var sites []string
	if sitesStr != "" {
		sites = strings.Split(sitesStr, ",")
	} else {
		// Если в системе пусто, бахаем дефолтный список, чтобы ничего не падало
		sites = []string{
			"https://google.com",
			"https://github.com",
			"https://yandex.ru",
		}
	}

	return &Config{
		Port:          port,
		CheckInterval: interval,
		FilePath:      filePath,
		Sites:         sites,
	}
}