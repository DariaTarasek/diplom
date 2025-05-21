package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetAppointments Получение списка всех записей
func (s *Store) GetAppointments(ctx context.Context) ([]model.Appointment, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointments").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка записей: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var appointment []model.Appointment
	err = s.db.SelectContext(dbCtx, &appointment, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка записей: %w", err)
	}

	return appointment, nil
}

// GetAppointmentByID Получение записи по id
func (s *Store) GetAppointmentByID(ctx context.Context, id model.AppointmentID) (model.Appointment, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointments").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Appointment{}, fmt.Errorf("не удалось сформировать запрос для получения записи по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var appointment model.Appointment
	err = s.db.GetContext(dbCtx, &appointment, query, args...)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("не удалось выполнить запрос для получения записи по id: %w", err)
	}

	return appointment, nil
}

// AddAppointment Добавление новой записи
func (s *Store) AddAppointment(ctx context.Context, appointment model.Appointment) (model.AppointmentID, error) {
	fields := map[string]any{
		"patient_id":   appointment.PatientID,
		"doctor_id":    appointment.DoctorID,
		"date":         appointment.Date,
		"time":         appointment.Time,
		"second_name":  appointment.PatientSecondName,
		"first_name":   appointment.PatientFirstName,
		"surname":      appointment.PatientSurname,
		"birth_date":   appointment.PatientBirthDate,
		"gender":       appointment.PatientGender,
		"phone_number": appointment.PatientPhoneNumber,
		"status":       appointment.Status,
		"created_at":   appointment.CreatedAt,
		"updated_at":   appointment.UpdatedAt,
	}
	query, args, err := s.builder.
		Insert("appointments").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return model.AppointmentID(0), fmt.Errorf("не удалось сформировать запрос для добавления новой записи: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var appointmentID model.AppointmentID
	err = s.db.QueryRowxContext(dbCtx, query, args...).Scan(&appointmentID)
	if err != nil {
		return model.AppointmentID(0), fmt.Errorf("не удалось выполнить запрос для добавления новой записи: %w", err)
	}
	return appointmentID, nil
}

// UpdateAppointment Изменение данных записи
func (s *Store) UpdateAppointment(ctx context.Context, id model.AppointmentID, appointment model.Appointment) error {
	fields := map[string]any{
		"patient_id":   appointment.PatientID,
		"doctor_id":    appointment.DoctorID,
		"date":         appointment.Date,
		"time":         appointment.Time,
		"second_name":  appointment.PatientSecondName,
		"first_name":   appointment.PatientFirstName,
		"surname":      appointment.PatientSurname,
		"birth_date":   appointment.PatientBirthDate,
		"gender":       appointment.PatientGender,
		"phone_number": appointment.PatientPhoneNumber,
		"status":       appointment.Status,
		"created_at":   appointment.CreatedAt,
		"updated_at":   appointment.UpdatedAt,
	}
	query, args, err := s.builder.
		Update("appointments").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных записи: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных записи: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных записи: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных записи: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные записи не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных записи: %w", err)
	}
	return nil
}

// DeleteAppointment Удаление записи по id
func (s *Store) DeleteAppointment(ctx context.Context, id model.AppointmentID) error {
	query, args, err := s.builder.
		Delete("appointments").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления записи: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления записи: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления записи: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления записи: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("запись не удалена")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления записи: %w", err)
	}
	return nil
}
