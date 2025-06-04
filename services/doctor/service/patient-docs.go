package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/doctor/model"
	storagepb "github.com/DariaTarasek/diplom/services/doctor/proto/storage"
	"strconv"
)

func (s *DoctorService) GetDocumentsInfo(ctx context.Context, patientID int) ([]model.DocumentInfo, error) {

	uid := strconv.Itoa(patientID)
	resp, err := s.StorageClient.Client.GetDocumentsByPatientID(ctx, &storagepb.GetDocumentsRequest{PatientId: uid})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить информацию о документах: %w", err)
	}

	docsInfo := make([]model.DocumentInfo, 0, len(resp.Documents))
	for _, item := range resp.Documents {
		docInfo := model.DocumentInfo{
			ID:          item.Id,
			FileName:    item.FileName,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.AsTime().Format("02.01.2006"),
		}
		docsInfo = append(docsInfo, docInfo)
	}
	return docsInfo, nil
}

func (s *DoctorService) DownloadDocument(ctx context.Context, documentID string) (model.DocumentFile, error) {

	// Запрашиваем документ по ID через storage-сервис
	resp, err := s.StorageClient.Client.DownloadDocument(ctx, &storagepb.DownloadDocumentRequest{DocumentId: documentID})
	if err != nil {
		return model.DocumentFile{}, fmt.Errorf("не удалось скачать документ: %w", err)
	}

	// Собираем и возвращаем структуру модели
	return model.DocumentFile{
		FileName:    resp.FileName,
		FileContent: resp.FileContent,
	}, nil
}
