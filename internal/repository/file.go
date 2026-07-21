package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"uptime-checker/internal/domain"
)

type FileRepository struct {
	filePath string
}

func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{filePath: filePath}
}

// Save сохраняет результат проверки в файл.
func (r *FileRepository) Save(result domain.CheckResult) {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[ОШИБКА ФАЙЛА] Не удалось открыть файл: %v\n", err)
		return
	}
	defer file.Close()

	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Printf("[ОШИБКА JSON] Не удалось закодировать данные: %v\n", err)
		return
	}

	_, err = file.Write(append(jsonData, '\n'))
	if err != nil {
		fmt.Printf("[ОШИБКА ЗАПИСИ] Не удалось записать в файл: %v\n", err)
	}
}

// ReadAll читает весь файл с диска.
func (r *FileRepository) ReadAll() ([]byte, error) {
	return os.ReadFile(r.filePath)
}