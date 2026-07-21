package main

import (
	"fmt"
	"net/http"
	"time"
	"uptime-checker/internal/handler"
	"uptime-checker/internal/repository"
	"uptime-checker/internal/service"
)

func main() {
	sites := []string{
		"https://google.com",
		"https://github.com",
		"https://yandex.ru",
		"https://httpbin.org",
	}

	// Инициализируем слои приложения (Dependency Injection)
	repo := repository.NewFileRepository("metrics.json")
	svc := service.NewMonitorService(repo, len(sites))
	h := handler.NewHTTPHandler(repo)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Запуск HTTP-сервера в фоне
	http.HandleFunc("/metrics", h.MetricsHandler)
	go func() {
		fmt.Println("[HTTP] Сервер метрик доступен на http://localhost:8080/metrics")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("[ОШИБКА HTTP] Не удалось запустить сервер: %v\n", err)
		}
	}()

	fmt.Println("--- Сервис запущен. Контроль каждые 10 сек. ---")
	
	// Первый запуск
	svc.RunCheckRound(sites)

	// Цикл мониторинга по таймеру
	for range ticker.C {
		fmt.Println("\n[Таймер сработал] Начинаем новый круг проверки. . .")
		svc.RunCheckRound(sites)
	}
}