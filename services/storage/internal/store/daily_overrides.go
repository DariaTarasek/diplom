package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
	"time"
)

// МЕТОДЫ ДЛЯ ТАБЛИЦЫ DOCTOR_DAILY_OVERRIDE

// GetDoctorsOverrides Получение списка переопределений расписания всех врачей
func (s *Store) GetDoctorsOverrides(ctx context.Context) ([]model.DoctorDailyOverride, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_daily_override").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка переопределений расписания врачей: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrides []model.DoctorDailyOverride
	err = s.db.SelectContext(dbCtx, &overrides, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка переопределений расписания: %w", err)
	}

	return overrides, nil
}

// GetDoctorsOverridesByDate Получение переопределения расписания врачей по дате
func (s *Store) GetDoctorsOverridesByDate(ctx context.Context, date time.Time) ([]model.DoctorDailyOverride, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_daily_override").
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения переопределения расписания врача по дате: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrides []model.DoctorDailyOverride
	err = s.db.SelectContext(dbCtx, &overrides, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения переопределения расписания врача по дате: %w", err)
	}

	return overrides, nil
}

// GetOverridesByDoctorID Получение переопределения расписания по врачу
func (s *Store) GetOverridesByDoctorID(ctx context.Context, id model.UserID) ([]model.DoctorDailyOverride, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_daily_override").
		Where(squirrel.Eq{"doctor_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения переопределения расписания по врачу: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrides []model.DoctorDailyOverride
	err = s.db.SelectContext(dbCtx, &overrides, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения переопределения расписания по врачу: %w", err)
	}

	return overrides, nil
}

// GetOverridesByDoctorAndDate Получение переопределения расписания по врачу и дате
func (s *Store) GetOverridesByDoctorAndDate(ctx context.Context, id model.UserID, date time.Time) ([]model.DoctorDailyOverride, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_daily_override").
		Where(squirrel.Eq{"doctor_id": id, "date": date}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения переопределения расписания по врачу и дате: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrides []model.DoctorDailyOverride
	err = s.db.SelectContext(dbCtx, &overrides, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения переопределения расписания по врачу и дате: %w", err)
	}

	return overrides, nil
}

// AddDoctorOverride Добавление нового переопределения расписания врача
func (s *Store) AddDoctorOverride(ctx context.Context, override model.DoctorDailyOverride) (int, error) {
	fields := map[string]any{
		"doctor_id":             override.DoctorID,
		"date":                  override.Date,
		"start_time":            override.StartTime,
		"end_time":              override.EndTime,
		"slot_duration_minutes": override.SlotDurationMinutes,
		"is_day_off":            override.IsDayOff,
	}
	query, args, err := s.builder.
		Insert("doctor_daily_override").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для добавления нового переопределения расписания врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrideID int
	err = s.db.QueryRowxContext(dbCtx, query, args...).Scan(&overrideID)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для добавления нового переопределения расписания врача: %w", err)
	}
	return overrideID, nil
}

// UpdateDoctorOverride Изменение данных переопределения расписания
func (s *Store) UpdateDoctorOverride(ctx context.Context, id int, override model.DoctorDailyOverride) error {
	fields := map[string]any{
		"date":                  override.Date,
		"start_time":            override.StartTime,
		"end_time":              override.EndTime,
		"slot_duration_minutes": override.SlotDurationMinutes,
		"is_day_off":            override.IsDayOff,
	}
	query, args, err := s.builder.
		Update("doctor_daily_override").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных переопределения расписания врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных переопределения расписания врача: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных переопределения расписания врача: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных переопределениия расписания врача: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные переопределения расписания врача не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных переопределения расписания врача: %w", err)
	}
	return nil
}

// DeleteDoctorOverride Удаление переопределения расписания врача
func (s *Store) DeleteDoctorOverride(ctx context.Context, id int) error {
	query, args, err := s.builder.
		Delete("doctor_daily_override").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления переопределения расписания врача: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления переопределения расписания врача: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления переопределения расписания врача: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления переопределения расписания врача: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("переопределение расписания врача не удалено")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления переопределения расписания врача: %w", err)
	}
	return nil
}

// МЕТОДЫ ДЛЯ ТАБЛИЦЫ CLINIC_DAILY_OVERRIDE

// GetClinicOverrides Получение списка переопределений расписания клиники
func (s *Store) GetClinicOverrides(ctx context.Context) ([]model.ClinicDailyOverride, error) {
	query, args, err := s.builder.
		Select("*").
		From("clinic_daily_override").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка переопределений расписания клиники: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrides []model.ClinicDailyOverride
	err = s.db.SelectContext(dbCtx, &overrides, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка переопределений расписания клиники: %w", err)
	}

	return overrides, nil
}

// GetClinicOverridesByDate Получение переопределения расписания врачей по дате
func (s *Store) GetClinicOverridesByDate(ctx context.Context, date time.Time) ([]model.ClinicDailyOverride, error) {
	query, args, err := s.builder.
		Select("*").
		From("clinic_daily_override").
		Where(squirrel.Eq{"date": date}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения переопределения расписания клиники по дате: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrides []model.ClinicDailyOverride
	err = s.db.SelectContext(dbCtx, &overrides, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения переопределения расписания клиники по дате: %w", err)
	}

	return overrides, nil
}

// AddClinicOverride Добавление нового переопределения расписания клиники
func (s *Store) AddClinicOverride(ctx context.Context, override model.ClinicDailyOverride) (int, error) {
	fields := map[string]any{
		"date":                  override.Date,
		"start_time":            override.StartTime,
		"end_time":              override.EndTime,
		"slot_duration_minutes": override.SlotDurationMinutes,
		"is_day_off":            override.IsDayOff,
	}
	query, args, err := s.builder.
		Insert("clinic_daily_override").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для добавления нового переопределения расписания клиники: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var overrideID int
	err = s.db.QueryRowxContext(dbCtx, query, args...).Scan(&overrideID)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для добавления нового переопределения расписания клиники: %w", err)
	}
	return overrideID, nil
}

// UpdateClinicOverride Изменение данных переопределения расписания
func (s *Store) UpdateClinicOverride(ctx context.Context, id int, override model.ClinicDailyOverride) error {
	fields := map[string]any{
		"date":                  override.Date,
		"start_time":            override.StartTime,
		"end_time":              override.EndTime,
		"slot_duration_minutes": override.SlotDurationMinutes,
		"is_day_off":            override.IsDayOff,
	}
	query, args, err := s.builder.
		Update("clinic_daily_override").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных переопределения расписания клиники: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных переопределения расписания клиники: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных переопределения расписания клиники: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных переопределениия расписания клиники: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные переопределения расписания клиники не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных переопределения расписания клиники: %w", err)
	}
	return nil
}

// DeleteClinicOverride Удаление переопределения расписания клиники
func (s *Store) DeleteClinicOverride(ctx context.Context, id int) error {
	query, args, err := s.builder.
		Delete("clinic_daily_override").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления переопределения расписания клиники: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления переопределения расписания клиники: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления переопределения расписания клиники: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления переопределения расписания клиники: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("переопределение расписания клиники не удалено")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления переопределения расписания клиники: %w", err)
	}
	return nil
}
