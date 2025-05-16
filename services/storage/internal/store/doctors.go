package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetDoctors Получение списка всех врачей
func (s *Store) GetDoctors(ctx context.Context) ([]model.Doctor, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctors").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка врачей: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var doctors []model.Doctor
	err = s.db.SelectContext(dbCtx, &doctors, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка врачей: %w", err)
	}

	return doctors, nil
}

// GetDoctorByID Получение врача по id
func (s *Store) GetDoctorByID(ctx context.Context, id model.UserID) (model.Doctor, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctors").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return model.Doctor{}, fmt.Errorf("не удалось сформировать запрос для получения врачa: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var doctor model.Doctor
	err = s.db.SelectContext(dbCtx, &doctor, query, args...)
	if err != nil {
		return model.Doctor{}, fmt.Errorf("не удалось выполнить запрос для получения врачa: %w", err)
	}

	return doctor, nil
}

// AddDoctor Добавление нового врача
func (s *Store) AddDoctor(ctx context.Context, doctor model.Doctor) error {
	fields := map[string]any{
		"user_id":      doctor.ID,
		"first_name":   doctor.FirstName,
		"second_name":  doctor.SecondName,
		"surname":      doctor.Surname,
		"phone_number": doctor.PhoneNumber,
		"email":        doctor.Email,
		"education":    doctor.Education,
		"experience":   doctor.Experience,
		"gender":       doctor.Gender,
	}
	query, args, err := s.builder.
		Insert("doctors").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления нового врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res, err := s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления нового врача: %w", err)
	}
	_, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("не удалось получить id добавленного врача: %w", err)
	}
	return nil
}

// UpdateDoctor Изменение данных врача
func (s *Store) UpdateDoctor(ctx context.Context, id model.UserID, doctor model.Doctor) error {
	fields := map[string]any{
		"first_name":   doctor.FirstName,
		"second_name":  doctor.SecondName,
		"surname":      doctor.Surname,
		"phone_number": doctor.PhoneNumber,
		"email":        doctor.Email,
		"education":    doctor.Education,
		"experience":   doctor.Experience,
		"gender":       doctor.Gender,
	}
	query, args, err := s.builder.
		Update("doctors").
		Where(squirrel.Eq{"user_id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных врача: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных врача: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных врача: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные врача не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных врача: %w", err)
	}
	return nil
}

// DeleteDoctor Удаление врача по id
func (s *Store) DeleteDoctor(ctx context.Context, id model.UserID) error {
	query, args, err := s.builder.
		Delete("doctors").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления врачa: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления врача: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления врачa: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления врача: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("доктор не удален")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления врача: %w", err)
	}
	return nil
}
