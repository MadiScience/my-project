package service

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"uptime-checker/internal/domain"
	"uptime-checker/internal/repository"
)

type MonitorService struct {
	repo        *repository.FileRepository
	resultsChan chan domain.CheckResult
}

func NewMonitorService(repo *repository.FileRepository, bufferSize int) *MonitorService {
	s := &MonitorService{
		repo:        repo,
		resultsChan: make(chan domain.CheckResult, bufferSize),
	}
	// Сразу запускаем фоновый воркер для обработки результатов
	go s.startWorker()
	return s
}

func (s *MonitorService) GetResultsChan() chan domain.CheckResult {
	return s.resultsChan
}

// startWorker слушает канал и отправляет данные в репозиторий
func (s *MonitorService) startWorker() {
	for result := range s.resultsChan {
		if result.IsAlive {
			fmt.Printf("[%s] [УСПЕХ] Сайт %s ответил за %v (Код: %d)\n", result.CheckedAt, result.URL, result.ResponseTime, result.StatusCode)
		} else {
			fmt.Printf("[%s] [ОШИБКА] Сайт %s НЕДОСТУПЕН! Время: %v\n", result.CheckedAt, result.URL, result.ResponseTime)
		}
		s.repo.Save(result)
	}
}

// RunCheckRound контролирует выполнение горутин текущего круга
func (s *MonitorService) RunCheckRound(sites []string) {
	var wg sync.WaitGroup

	for _, siteURL := range sites {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			s.resultsChan <- s.checkSite(url)
		}(siteURL)
	}
	wg.Wait()
	fmt.Println("[Инфо] Все горутины текущего раунда успешно отработали.")
}

func (s *MonitorService) checkSite(url string) domain.CheckResult {
	startTime := time.Now()
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	duration := time.Since(startTime)
	checkedAt := time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		return domain.CheckResult{
			URL:          url,
			StatusCode:   0,
			ResponseTime: duration,
			IsAlive      false,
			CheckedAt:    checkedAt,
		}
	}
	defer resp.Body.Close()

	return domain.CheckResult{
		URL:          url,
		StatusCode:   resp.StatusCode,
		ResponseTime: duration,
		IsAlive:      resp.StatusCode >= 200 && resp.StatusCode < 300,
		CheckedAt:    checkedAt,
	}
}