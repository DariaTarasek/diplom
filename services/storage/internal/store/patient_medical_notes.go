package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetPatientMedicalNotes Получение списка всей мед. информации о пациенте
func (s *Store) GetPatientMedicalNotes(ctx context.Context, id model.UserID) ([]model.PatientMedicalNote, error) {
	query, args, err := s.builder.
		Select("*").
		From("patient_medical_notes").
		Where(squirrel.Eq{"patient_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения мед. инф. пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var notes []model.PatientMedicalNote
	err = s.db.SelectContext(dbCtx, &notes, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения мед. инф. пациента: %w", err)
	}

	return notes, nil
}

// AddPatientMedicalNote Добавление новой мед. инф. пациента
func (s *Store) AddPatientMedicalNote(ctx context.Context, note model.PatientMedicalNote) error {
	fields := map[string]any{
		"patient_id": note.PatientID,
		"type":       note.Type,
		"title":      note.Title,
		"created_at": note.CreatedAt,
	}
	query, args, err := s.builder.
		Insert("patient_medical_notes").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления новой новой мед. инф. пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления новой мед. инф. пациента: %w", err)
	}

	return nil
}

// UpdatePatientMedicalNotes Изменение мед. инф. пациента
func (s *Store) UpdatePatientMedicalNotes(ctx context.Context, id int, note model.PatientMedicalNote) error {
	fields := map[string]any{
		"title":      note.Title,
		"created_at": note.CreatedAt,
	}
	query, args, err := s.builder.
		Update("patient_medical_notes").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления мед. инф. пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения мед. инф. пациента: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения мед. инф. пациента: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения мед. инф. пациента: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("мед. инф. пациента не измененa")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения мед. инф. пациента: %w", err)
	}
	return nil
}
