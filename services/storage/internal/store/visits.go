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
		"complaints":     visit.Complaints,
		"treatment_plan": visit.TreatmentPlan,
		"created_at":     visit.CreatedAt,
		"updated_at":     visit.UpdatedAt,
	}
	query, args, err := s.builder.
		Insert("appointment_visits").
		SetMap(fields).
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
