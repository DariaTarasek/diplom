package storagefs

import (
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"os"
	"path/filepath"
)

type FileStorage struct {
	BasePath string
}

func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{BasePath: basePath}
}

func (s *FileStorage) SaveFile(patientID, filename string, data []byte) (string, error) {
	dir := filepath.Join(s.BasePath, patientID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("не удалось записать файл: %w", err)
	}

	return fullPath, nil
}

func (s *FileStorage) GetDocumentFile(doc *model.PatientDocument) ([]byte, string, error) {
	data, err := os.ReadFile(doc.StoragePath)
	if err != nil {
		return nil, "", fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	return data, doc.FileName, nil
}
