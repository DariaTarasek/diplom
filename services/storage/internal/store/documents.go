package store

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (s *Store) SaveDocument(ctx context.Context, doc *model.PatientDocument) error {
	fields := map[string]any{
		"id":           doc.ID,
		"patient_id":   doc.PatientID,
		"file_name":    doc.FileName,
		"modality":     doc.Modality,
		"study_date":   doc.StudyDate,
		"description":  doc.Description,
		"storage_path": doc.StoragePath,
		"preview_path": doc.PreviewPath,
		"created_at":   doc.CreatedAt,
	}
	query, args, err := s.builder.
		Insert("patient_documents").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось создать запрос для сохранения данных о документе: %w", err)
	}
	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для сохранения данных о документе: %w", err)
	}
	return nil
}

func (s *Store) GetDocumentsByPatient(ctx context.Context, patientID model.UserID) ([]model.PatientDocument, error) {
	query, args, err := s.builder.
		Select("*").
		From("patient_documents").
		Where(squirrel.Eq{"patient_id": patientID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос для получения информации о документе: %w", err)
	}

	var docs []model.PatientDocument
	if err := s.db.SelectContext(ctx, &docs, query, args...); err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения информации о документе: %w", err)
	}
	return docs, nil
}

func (s *Store) GetDocumentByID(ctx context.Context, id uuid.UUID) (*model.PatientDocument, error) {
	query, args, err := s.builder.
		Select("*").
		From("patient_documents").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения информации о документе: %w", err)
	}

	var doc model.PatientDocument
	if err := s.db.GetContext(ctx, &doc, query, args...); err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения информации о документе: %w", err)
	}
	return &doc, nil
}
