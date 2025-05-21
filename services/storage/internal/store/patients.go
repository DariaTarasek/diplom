package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetPatients Получение списка всех пациентов
func (s *Store) GetPatients(ctx context.Context) ([]model.Patient, error) {
	query, args, err := s.builder.
		Select("*").
		From("patients").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка пациентов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var patients []model.Patient
	err = s.db.SelectContext(dbCtx, &patients, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка пациентов: %w", err)
	}

	return patients, nil
}

// GetPatientByID Получение пациента по ID
func (s *Store) GetPatientByID(ctx context.Context, id model.UserID) (model.Patient, error) {
	query, args, err := s.builder.
		Select("*").
		From("patients").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return model.Patient{}, fmt.Errorf("не удалось сформировать запрос для получения пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var patient model.Patient
	err = s.db.GetContext(dbCtx, &patient, query, args...)
	if err != nil {
		return model.Patient{}, fmt.Errorf("не удалось выполнить запрос для получения пациента: %w", err)
	}

	return patient, nil
}

// AddPatient Добавление нового пациента
func (s *Store) AddPatient(ctx context.Context, patient model.Patient) error {
	fields := map[string]any{
		"user_id":      patient.ID,
		"first_name":   patient.FirstName,
		"second_name":  patient.SecondName,
		"surname":      patient.Surname,
		"phone_number": patient.PhoneNumber,
		"email":        patient.Email,
		"birth_date":   patient.BirthDate,
		"gender":       patient.Gender,
	}
	query, args, err := s.builder.
		Insert("patients").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления нового пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления нового пациента: %w", err)
	}
	return nil
}

// UpdatePatient Изменение данных пациента
func (s *Store) UpdatePatient(ctx context.Context, id model.UserID, patient model.Patient) error {
	fields := map[string]any{
		"first_name":   patient.FirstName,
		"second_name":  patient.SecondName,
		"surname":      patient.Surname,
		"phone_number": patient.PhoneNumber,
		"email":        patient.Email,
		"birth_date":   patient.BirthDate,
		"gender":       patient.Gender,
	}
	query, args, err := s.builder.
		Update("patients").
		Where(squirrel.Eq{"user_id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных пациента: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных пациента: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных пациента: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные пациента не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных пациента: %w", err)
	}
	return nil
}

// DeletePatient Удаление пациента по id
func (s *Store) DeletePatient(ctx context.Context, id model.UserID) error {
	query, args, err := s.builder.
		Delete("patients").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления пациента: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления пациента: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления пациента: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пациент не удален")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления пациента: %w", err)
	}
	return nil
}
