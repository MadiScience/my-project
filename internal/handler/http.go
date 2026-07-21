package handler

import (
	"net/http"
	"os"
	"uptime-checker/internal/repository"
)

type HTTPHandler struct {
	repo *repository.FileRepository
}

func NewHTTPHandler(repo *repository.FileRepository) *HTTPHandler {
	return &HTTPHandler{repo: repo}
}

func (h *HTTPHandler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	fileBytes, err := h.repo.ReadAll()
	if err != nil {
		if os.IsNotExist(err) {
			w.Write([]byte("[]"))
			return
		}
		http.Error(w, `{"error": "не удалось прочитать файл метрик"}`, http.StatusInternalServerError)
		return
	}

	// Твоя фирменная побайтовая сборка JSON-массива
	var finalJSON []byte
	finalJSON = append(finalJSON, '[')
	for i, b := range fileBytes {
		if b == '\n' {
			if i < len(fileBytes)-1 {
				finalJSON = append(finalJSON, ',')
			}
		} else {
			finalJSON = append(finalJSON, b)
		}
	}
	finalJSON = append(finalJSON, ']')

	w.Write(finalJSON)
}