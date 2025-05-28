package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// МЕТОДЫ ДЛЯ ТАБЛИЦЫ DOCTOR_WEEKLY_SCHEDULE

// GetDoctorsSchedules Получение расписания всех врачей
func (s *Store) GetDoctorsSchedules(ctx context.Context) ([]model.DoctorSchedule, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_weekly_schedule").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения расписания врачей: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var schedules []model.DoctorSchedule
	err = s.db.SelectContext(dbCtx, &schedules, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения расписания врачей: %w", err)
	}

	return schedules, nil
}

// GetDoctorsSchedulesByWeekday Получение расписания врачей по дню недели
func (s *Store) GetDoctorsSchedulesByWeekday(ctx context.Context, weekday int) ([]model.DoctorSchedule, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_weekly_schedule").
		Where(squirrel.Eq{"weekday": weekday}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения расписания врачей по дню недели: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var schedules []model.DoctorSchedule
	err = s.db.SelectContext(dbCtx, &schedules, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения расписания врачей по дню недели: %w", err)
	}

	return schedules, nil
}

// GetScheduleByDoctorID Получение расписания врача
func (s *Store) GetScheduleByDoctorID(ctx context.Context, id model.UserID) ([]model.DoctorSchedule, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_weekly_schedule").
		Where(squirrel.Eq{"doctor_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения расписания врачa: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var schedules []model.DoctorSchedule
	err = s.db.SelectContext(dbCtx, &schedules, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения расписания врачa: %w", err)
	}

	return schedules, nil
}

// GetScheduleByDoctorAndWeekday Получение расписания врачa по дню недели
func (s *Store) GetScheduleByDoctorAndWeekday(ctx context.Context, id model.UserID, weekday int) ([]model.DoctorSchedule, error) {
	query, args, err := s.builder.
		Select("*").
		From("doctor_weekly_schedule").
		Where(squirrel.Eq{"doctor_id": id, "weekday": weekday}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения расписания врачa по дню недели: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var schedule []model.DoctorSchedule
	err = s.db.SelectContext(dbCtx, &schedule, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения расписания врача по дню недели: %w", err)
	}

	return schedule, nil
}

// AddDoctorSchedule Добавление расписания врача на всю неделю
func (s *Store) AddDoctorSchedule(ctx context.Context, schedules []model.DoctorSchedule) ([]int, error) {
	if len(schedules) == 0 {
		return nil, fmt.Errorf("список расписаний пуст")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию для добавления расписания врача: %w", err)
	}

	defer tx.Rollback()

	builder := s.builder.Insert("doctor_weekly_schedule").
		Columns("doctor_id", "weekday", "start_time", "end_time", "slot_duration_minutes", "is_day_off")

	for _, schedule := range schedules {
		builder = builder.Values(
			schedule.DoctorID,
			schedule.Weekday,
			schedule.StartTime,
			schedule.EndTime,
			schedule.SlotDurationMinutes,
			schedule.IsDayOff,
		)
	}
	builder = builder.Suffix("RETURNING id")

	query, args, buildErr := builder.ToSql()
	if buildErr != nil {
		err = fmt.Errorf("не удалось сформировать запрос для добавления расписания врача: %w", buildErr)
		return nil, err
	}

	rows, queryErr := tx.QueryxContext(ctx, query, args...)
	if queryErr != nil {
		err = fmt.Errorf("ошибка выполнения запроса добавления расписания врача: %w", queryErr)
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if scanErr := rows.Scan(&id); scanErr != nil {
			err = fmt.Errorf("ошибка чтения ID добавленного расписания врача: %w", scanErr)
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию для добавления расписания врача: %w", err)
	}

	return ids, nil
}

// UpdateDoctorSchedule Обновление постоянного расписания врача
func (s *Store) UpdateDoctorSchedule(ctx context.Context, schedules []model.DoctorSchedule) error {
	if len(schedules) == 0 {
		return fmt.Errorf("пустое расписание — обновление невозможно")
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для обновления расписания врача: %w", err)
	}

	for _, schedule := range schedules {
		fields := map[string]any{
			"start_time":            schedule.StartTime,
			"end_time":              schedule.EndTime,
			"slot_duration_minutes": schedule.SlotDurationMinutes,
			"is_day_off":            schedule.IsDayOff,
		}

		// Попытка обновить
		updateQuery, updateArgs, err := s.builder.
			Update("doctor_weekly_schedule").
			SetMap(fields).
			Where(squirrel.Eq{
				"doctor_id": schedule.DoctorID,
				"weekday":   schedule.Weekday,
			}).ToSql()

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось сформировать запрос на обновление расписания для дня %v: %w", schedule.Weekday, err)
		}

		res, err := tx.ExecContext(dbCtx, updateQuery, updateArgs...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось выполнить обновление расписания для дня %v: %w", schedule.Weekday, err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось получить количество затронутых строк при обновлении дня %v: %w", schedule.Weekday, err)
		}

		// Если ничего не обновилось — вставляем новую запись
		if rowsAffected == 0 {
			insertQuery, insertArgs, err := s.builder.
				Insert("doctor_weekly_schedule").
				Columns("doctor_id", "weekday", "start_time", "end_time", "slot_duration_minutes", "is_day_off").
				Values(
					schedule.DoctorID,
					schedule.Weekday,
					schedule.StartTime,
					schedule.EndTime,
					schedule.SlotDurationMinutes,
					schedule.IsDayOff,
				).ToSql()

			if err != nil {
				tx.Rollback()
				return fmt.Errorf("не удалось сформировать запрос на вставку расписания для дня %v: %w", schedule.Weekday, err)
			}

			_, err = tx.ExecContext(dbCtx, insertQuery, insertArgs...)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("не удалось вставить расписание для дня %v: %w", schedule.Weekday, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию обновления расписания: %w", err)
	}
	return nil
}

// МЕТОДЫ ДЛЯ ТАБЛИЦЫ CLINIC_WEEKLY_SCHEDULE

// GetClinicSchedule Получение расписания клиники
func (s *Store) GetClinicSchedule(ctx context.Context) ([]model.ClinicSchedule, error) {
	query, args, err := s.builder.
		Select("*").
		From("clinic_weekly_schedule").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения расписания клиники: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var schedule []model.ClinicSchedule
	err = s.db.SelectContext(dbCtx, &schedule, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения расписания клиники: %w", err)
	}

	return schedule, nil
}

// GetClinicScheduleByWeekday Получение расписания клиники по дню недели
func (s *Store) GetClinicScheduleByWeekday(ctx context.Context, weekday int) ([]model.ClinicSchedule, error) {
	query, args, err := s.builder.
		Select("*").
		From("clinic_weekly_schedule").
		Where(squirrel.Eq{"weekday": weekday}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения расписания клиники по дню недели: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var schedule []model.ClinicSchedule
	err = s.db.SelectContext(dbCtx, &schedule, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения расписания клиники по дню недели: %w", err)
	}

	return schedule, nil
}

// AddClinicSchedule Добавление расписания клиники на неделю
func (s *Store) AddClinicSchedule(ctx context.Context, schedules []model.ClinicSchedule) ([]int, error) {
	if len(schedules) == 0 {
		return nil, fmt.Errorf("список расписаний пуст")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию для добавления расписания клиники: %w", err)
	}

	defer tx.Rollback()

	builder := s.builder.Insert("clinic_weekly_schedule").
		Columns("weekday", "start_time", "end_time", "slot_duration_minutes", "is_day_off")

	for _, schedule := range schedules {
		builder = builder.Values(
			schedule.Weekday,
			schedule.StartTime,
			schedule.EndTime,
			schedule.SlotDurationMinutes,
			schedule.IsDayOff,
		)
	}
	builder = builder.Suffix("RETURNING id")

	query, args, buildErr := builder.ToSql()
	if buildErr != nil {
		err = fmt.Errorf("не удалось сформировать запрос для добавления расписания клиники: %w", buildErr)
		return nil, err
	}

	rows, queryErr := tx.QueryxContext(ctx, query, args...)
	if queryErr != nil {
		err = fmt.Errorf("ошибка выполнения запроса добавления расписания клиники: %w", queryErr)
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if scanErr := rows.Scan(&id); scanErr != nil {
			err = fmt.Errorf("ошибка чтения ID добавленного расписания клиники: %w", scanErr)
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию для добавления расписания клиники: %w", err)
	}
	return ids, nil
}

// UpdateClinicSchedule Изменение расписания клиники на неделю
func (s *Store) UpdateClinicSchedule(ctx context.Context, schedules []model.ClinicSchedule) error {
	if len(schedules) == 0 {
		return fmt.Errorf("пустое расписание — обновление невозможно")
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для обновления расписания клиники: %w", err)
	}

	for _, schedule := range schedules {
		fields := map[string]any{
			"start_time":            schedule.StartTime,
			"end_time":              schedule.EndTime,
			"slot_duration_minutes": schedule.SlotDurationMinutes,
			"is_day_off":            schedule.IsDayOff,
		}

		query, args, err := s.builder.
			Update("clinic_weekly_schedule").
			Where(squirrel.Eq{
				"weekday": schedule.Weekday,
			}).
			SetMap(fields).
			ToSql()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось сформировать запрос на обновление расписания для дня %v: %w", schedule.Weekday, err)
		}

		res, err := tx.ExecContext(dbCtx, query, args...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось выполнить обновление расписания для дня %v: %w", schedule.Weekday, err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось получить количество затронутых строк при обновлении дня %v: %w", schedule.Weekday, err)
		}

		if rowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("расписание на %v не было обновлено", schedule.Weekday)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию обновления расписания: %w", err)
	}

	return nil
}
