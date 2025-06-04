package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

func (s *Store) GetVisits(ctx context.Context) ([]model.Visit, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_visits").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка приемов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var visits []model.Visit
	err = s.db.SelectContext(dbCtx, &visits, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка приемов: %w", err)
	}

	return visits, nil
}

func (s *Store) GetVisitsByPatientID(ctx context.Context, patientID model.UserID) ([]model.Visit, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_visits").
		Where(squirrel.Eq{"patient_id": patientID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка приемов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var visits []model.Visit
	err = s.db.SelectContext(dbCtx, &visits, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка приемов: %w", err)
	}

	return visits, nil
}

func (s *Store) GetVisitByID(ctx context.Context, id model.VisitID) (model.Visit, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_visits").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Visit{}, fmt.Errorf("не удалось сформировать запрос для получения приема по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var visit model.Visit
	err = s.db.GetContext(dbCtx, &visit, query, args...)
	if err != nil {
		return model.Visit{}, fmt.Errorf("не удалось выполнить запрос для получения приема по id: %w", err)
	}

	return visit, nil
}

func (s *Store) GetVisitsByPatientId(ctx context.Context, id model.UserID) ([]model.Visit, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_visits").
		Where(squirrel.Eq{"patient_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения визитов пациента: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var visits []model.Visit
	err = s.db.SelectContext(dbCtx, &visits, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения визитов пациента: %w", err)
	}

	return visits, nil
}

func (s *Store) GetVisitByAppointmentID(ctx context.Context, id model.AppointmentID) (model.Visit, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_visits").
		Where(squirrel.Eq{"appointment_id": id}).
		ToSql()
	if err != nil {
		return model.Visit{}, fmt.Errorf("не удалось сформировать запрос для получения приема по id записи: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var visit model.Visit
	err = s.db.GetContext(dbCtx, &visit, query, args...)
	if err != nil {
		return model.Visit{}, fmt.Errorf("не удалось выполнить запрос для получения приема по id записи: %w", err)
	}

	return visit, nil
}

func (s *Store) AddVisit(ctx context.Context, visit model.Visit) (model.VisitID, error) {
	fields := map[string]any{
		"appointment_id": visit.AppointmentID,
		"patient_id":     visit.PatientID,
		"doctor_id":      visit.DoctorID,
		"complaints":     visit.Complaints,
		"treatment_plan": visit.TreatmentPlan,
		"created_at":     visit.CreatedAt,
		"updated_at":     visit.UpdatedAt,
	}
	query, args, err := s.builder.
		Insert("appointment_visits").
		SetMap(fields).Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для добавления приема: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var id model.VisitID
	err = s.db.QueryRowContext(dbCtx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для добавления приема: %w", err)
	}

	return id, nil
}

func (s *Store) UpdateVisit(ctx context.Context, id model.VisitID, visit model.Visit) error {
	fields := map[string]any{
		"appointment_id": visit.AppointmentID,
		"patient_id":     visit.PatientID,
		"doctor_id":      visit.DoctorID,
		"complaints":     visit.Complaints,
		"treatment_plan": visit.TreatmentPlan,
		"created_at":     visit.CreatedAt,
		"updated_at":     visit.UpdatedAt,
	}
	query, args, err := s.builder.
		Update("appointment_visits").
		SetMap(fields).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для изменения приема: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения приема: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения приема: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество обновленных приемов: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("прием не был обновлен")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию изменения приема: %w", err)
	}

	return nil
}

// AddVisitPayment Добавление стоимости приема
func (s *Store) AddVisitPayment(ctx context.Context, payment model.VisitPayment) error {
	fields := map[string]any{
		"visit_id": payment.VisitID,
		"price":    payment.Price,
		"status":   payment.Status,
	}
	query, args, err := s.builder.
		Insert("visit_payments").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления платежа: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления платежа: %w", err)
	}

	return nil
}

// UpdateVisitPayment Обновление стоимости приема
func (s *Store) UpdateVisitPayment(ctx context.Context, visitID model.VisitID, payment model.VisitPayment) error {
	fields := map[string]any{
		"price":  payment.Price,
		"status": payment.Status,
	}

	query, args, err := s.builder.
		Update("visit_payments").
		SetMap(fields).
		Where(squirrel.Eq{"visit_id": visitID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для изменения платежа: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения платежа: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения платежа: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество обновленных платежей: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("платеж не был обновлен")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию изменения платежа: %w", err)
	}

	return nil
}

func (s *Store) CalculateVisitTotal(ctx context.Context, visitID model.VisitID) (int, error) {
	var totalServices int
	queryServices, args, err := s.builder.
		Select("COALESCE(SUM(s.price * vs.quantity), 0)").
		From("appointment_services vs").
		Join("services s ON vs.service_id = s.id").
		Where(squirrel.Eq{"vs.visit_id": visitID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ошибка построения запроса услуг: %w", err)
	}
	err = s.db.QueryRowContext(ctx, queryServices, args...).Scan(&totalServices)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчёта суммы услуг: %w", err)
	}

	var totalMaterials int
	queryMaterials, args, err := s.builder.
		Select("COALESCE(SUM(m.price * vm.quantity_used), 0)").
		From("appointment_materials vm").
		Join("materials m ON vm.material_id = m.id").
		Where(squirrel.Eq{"vm.visit_id": visitID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ошибка построения запроса материалов: %w", err)
	}
	err = s.db.QueryRowContext(ctx, queryMaterials, args...).Scan(&totalMaterials)
	if err != nil {
		return 0, fmt.Errorf("ошибка подсчёта суммы материалов: %w", err)
	}

	return totalServices + totalMaterials, nil
}

func (s *Store) AddOrUpdateVisitPayment(ctx context.Context, payment model.VisitPayment) error {
	// Проверка существования записи
	existsQuery := `SELECT EXISTS(SELECT 1 FROM visit_payments WHERE visit_id = $1)`
	var exists bool
	err := s.db.QueryRowContext(ctx, existsQuery, payment.VisitID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования платежа: %w", err)
	}

	if exists {
		// Обновляем запись
		updateQuery, args, err := s.builder.
			Update("visit_payments").
			Set("price", payment.Price).
			Set("status", payment.Status).
			Where(squirrel.Eq{"visit_id": payment.VisitID}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("ошибка построения запроса обновления платежа: %w", err)
		}

		_, err = s.db.ExecContext(ctx, updateQuery, args...)
		if err != nil {
			return fmt.Errorf("ошибка выполнения запроса обновления платежа: %w", err)
		}
	} else {
		// Вставляем новую запись
		insertQuery, args, err := s.builder.
			Insert("visit_payments").
			Columns("visit_id", "price", "status").
			Values(payment.VisitID, payment.Price, payment.Status).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("ошибка построения запроса добавления платежа: %w", err)
		}

		_, err = s.db.ExecContext(ctx, insertQuery, args...)
		if err != nil {
			return fmt.Errorf("ошибка выполнения запроса добавления платежа: %w", err)
		}
	}

	return nil
}

func (s *Store) GetUnconfirmedVisitsPayments(ctx context.Context) ([]model.VisitPayment, error) {
	query, args, err := s.builder.
		Select("*").
		From("visit_payments").
		Where(squirrel.Eq{"status": "unconfirmed"}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var payments []model.VisitPayment
	if err := s.db.SelectContext(ctx, &payments, query, args...); err != nil {
		return nil, err
	}

	return payments, nil
}
