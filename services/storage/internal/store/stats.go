package store

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetTotalPatients Общее количество уникальных пациентов
func (s *Store) GetTotalPatients(ctx context.Context) (int, error) {
	query, args, err := s.builder.
		Select("COUNT(*)").
		From("patients").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для подсчета пациентов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var count int
	if err := s.db.GetContext(dbCtx, &count, query, args...); err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для подсчета пациентов: %w", err)
	}

	return count, nil
}

// GetTotalVisits Общее количество завершенных приёмов
func (s *Store) GetTotalVisits(ctx context.Context) (int, error) {
	query, args, err := s.builder.
		Select("COUNT(*)").
		From("appointment_visits").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для подсчета приёмов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var count int
	if err := s.db.GetContext(dbCtx, &count, query, args...); err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для подсчета приёмов: %w", err)
	}

	return count, nil
}

// GetTopServices возвращает топ-3 самых популярных услуг
func (s *Store) GetTopServices(ctx context.Context) ([]model.Top3Services, error) {
	query, args, err := s.builder.
		Select("s.name", "SUM(aps.quantity) AS usage_count").
		From("appointment_services aps").
		Join("services s ON aps.service_id = s.id").
		GroupBy("s.name").
		OrderBy("usage_count DESC").
		Limit(3).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для топ-услуг: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var stats []model.Top3Services
	if err := s.db.SelectContext(dbCtx, &stats, query, args...); err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для топ-услуг: %w", err)
	}

	return stats, nil
}

// GetDoctorWeeklyAverages Среднее количество приёмов в неделю у врача за всё время
func (s *Store) GetDoctorWeeklyAverages(ctx context.Context) ([]model.DoctorAvgVisit, error) {
	query, args, err := s.builder.
		Select("doctor_id", "ROUND(COUNT(*)::decimal / COUNT(DISTINCT DATE_TRUNC('week', created_at)), 2) AS avg_weekly_visits").
		From("appointment_visits").
		GroupBy("doctor_id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var stats []model.DoctorAvgVisit
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.SelectContext(dbCtx, &stats, query, args...); err != nil {
		return nil, fmt.Errorf("не удалось получить данные: %w", err)
	}
	return stats, nil
}

// GetDoctorAvgCheck Средний чек врача
func (s *Store) GetDoctorAvgCheck(ctx context.Context) ([]model.DoctorCheckStat, error) {
	query, args, err := s.builder.
		Select("av.doctor_id", "ROUND(AVG(vp.price), 2) AS avg_check").
		From("appointment_visits av").
		Join("visit_payments vp ON vp.visit_id = av.id").
		Where("vp.status = 'confirmed'").
		GroupBy("av.doctor_id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var stats []model.DoctorCheckStat
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.SelectContext(dbCtx, &stats, query, args...); err != nil {
		return nil, fmt.Errorf("не удалось получить данные: %w", err)
	}
	return stats, nil
}

// GetDoctorUniquePatients Количество уникальных пациентов врача
func (s *Store) GetDoctorUniquePatients(ctx context.Context) ([]model.DoctorUniquePatients, error) {
	query, args, err := s.builder.
		Select("doctor_id", "COUNT(DISTINCT patient_id) AS unique_patients").
		From("appointment_visits").
		GroupBy("doctor_id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var stats []model.DoctorUniquePatients
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.SelectContext(dbCtx, &stats, query, args...); err != nil {
		return nil, fmt.Errorf("не удалось получить данные: %w", err)
	}
	return stats, nil
}

// GetNewPatientsThisMonth Новые пациенты за текущий месяц
func (s *Store) GetNewPatientsThisMonth(ctx context.Context) (int, error) {
	subQuery := s.builder.
		Select("patient_id", "MIN(created_at)::date AS first_appointment_date").
		From("appointments").
		Where(squirrel.Eq{"status": "completed"}).
		GroupBy("patient_id")

	sqlSub, argsSub, err := subQuery.ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать подзапрос: %w", err)
	}

	query := "SELECT COUNT(*) FROM (" + sqlSub + ") AS first_appointments WHERE DATE_TRUNC('month', first_appointment_date) = DATE_TRUNC('month', CURRENT_DATE)"
	var count int
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.GetContext(dbCtx, &count, query, argsSub...); err != nil {
		return 0, fmt.Errorf("не удалось выполнить подзапрос: %w", err)
	}
	return count, nil
}

// GetAvgVisitsPerPatient Среднее количество визитов на пациента
func (s *Store) GetAvgVisitsPerPatient(ctx context.Context) (float64, error) {
	query, args, err := s.builder.
		Select("ROUND(COUNT(*)::decimal / COUNT(DISTINCT patient_id), 2) AS avg_visits").
		From("appointment_visits").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var avg float64
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.GetContext(dbCtx, &avg, query, args...); err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	return avg, nil
}

// GetAgeDistribution Возрастное распределение пациентов в процентах
func (s *Store) GetAgeDistribution(ctx context.Context) ([]model.AgeGroupStat, error) {
	// Подзапрос: сгруппированные возрастные данные
	ageCounts := s.builder.
		Select(
			"(EXTRACT(YEAR FROM AGE(birth_date)) / 10)::int * 10 AS age_group",
			"COUNT(*) AS count",
		).
		From("patients").
		GroupBy("age_group")

	ageCountsSQL, ageCountsArgs, err := ageCounts.ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать подзапрос ageCounts: %w", err)
	}

	finalQuery := fmt.Sprintf(`
		WITH age_counts AS (%s),
		     total AS (SELECT SUM(count) AS total_count FROM age_counts)
		SELECT 
			age_group,
			ROUND(count::decimal / total_count * 100, 2) AS percent
		FROM age_counts, total
		ORDER BY age_group
	`, ageCountsSQL)

	var stats []model.AgeGroupStat
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if err := s.db.SelectContext(dbCtx, &stats, finalQuery, ageCountsArgs...); err != nil {
		return nil, fmt.Errorf("не удалось получить данные: %w", err)
	}

	return stats, nil
}

// GetTotalIncome Общая выручка за всё время
func (s *Store) GetTotalIncome(ctx context.Context) (float64, error) {
	query, args, err := s.builder.
		Select("SUM(price) AS total_income").
		From("visit_payments").
		Where("status = 'confirmed'").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var income float64
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.GetContext(dbCtx, &income, query, args...); err != nil {
		return 0, fmt.Errorf("не удалось получить общую выручку: %w", err)
	}
	return income, nil
}

// GetMonthlyIncome Выручка за текущий месяц
func (s *Store) GetMonthlyIncome(ctx context.Context) (float64, error) {
	query, args, err := s.builder.
		Select("SUM(vp.price) AS monthly_income").
		From("visit_payments vp").
		Join("appointment_visits av ON vp.visit_id = av.id").
		Where("vp.status = 'confirmed'").
		Where("DATE_TRUNC('month', av.created_at) = DATE_TRUNC('month', CURRENT_DATE)").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var income float64
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.GetContext(dbCtx, &income, query, args...); err != nil {
		return 0, fmt.Errorf("не удалось получить месячную выручку: %w", err)
	}
	return income, nil
}

// GetClinicAverageCheck Средний чек по клинике
func (s *Store) GetClinicAverageCheck(ctx context.Context) (float64, error) {
	query, args, err := s.builder.
		Select("ROUND(AVG(price), 2) AS avg_check").
		From("visit_payments").
		Where("status = 'confirmed'").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос: %w", err)
	}
	var avg float64
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	if err := s.db.GetContext(dbCtx, &avg, query, args...); err != nil {
		return 0, fmt.Errorf("не удалось получить средний чек клиники: %w", err)
	}
	return avg, nil
}
