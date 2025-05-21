package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetDoctorsBySpecializationID Получение всех врачей выбранной специализации
func (s *Store) GetDoctorsBySpecializationID(ctx context.Context, id model.SpecID) ([]model.DoctorSpecialization, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_specializations").
		Where(squirrel.Eq{"specialization_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения врачей по специализации: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var specs []model.DoctorSpecialization
	err = s.db.SelectContext(dbCtx, &specs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения врачей по специализации: %w", err)
	}

	return specs, nil
}

// GetDoctorSpecializations Получение списка всех специализаций врача
func (s *Store) GetDoctorSpecializations(ctx context.Context, id model.UserID) ([]model.DoctorSpecialization, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_specializations").
		Where(squirrel.Eq{"doctor_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения специализаций врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var specs []model.DoctorSpecialization
	err = s.db.SelectContext(dbCtx, &specs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения специализаций врача: %w", err)
	}

	return specs, nil
}

// AddDoctorSpecialization Добавление новой специализации врачa
func (s *Store) AddDoctorSpecialization(ctx context.Context, spec model.DoctorSpecialization) error {
	fields := map[string]any{
		"doctor_id":         spec.DoctorID,
		"specialization_id": spec.SpecializationID,
	}
	query, args, err := s.builder.
		Insert("doctor_specializations").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления новой специализации врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления новой специализации врача: %w", err)
	}
	return nil
}

// DeleteDoctorSpecialization Удаление специализации у врача
func (s *Store) DeleteDoctorSpecialization(ctx context.Context, docID model.UserID, specID model.SpecID) error {
	query, args, err := s.builder.
		Delete("doctor_specializations").
		Where(squirrel.Eq{"doctor_id": docID, "specialization_id": specID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления специализации врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления специализации врача: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления специализации врача: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления специализации врача: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("специализация у врача не удалена")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления специализации врача: %w", err)
	}
	return nil
}
