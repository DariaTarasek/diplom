package store

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetAppointmentServices Получение списка всех оказанных во время приема услуг
func (s *Store) GetAppointmentServices(ctx context.Context, id model.VisitID) ([]model.AppointmentService, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_services").
		Where(squirrel.Eq{"visit_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения оказанных во время приема услуг: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var appServices []model.AppointmentService
	err = s.db.SelectContext(dbCtx, &appServices, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения оказанных во время приема услуг: %w", err)
	}

	return appServices, nil
}

// AddAppointmentServices Добавление записей о услугах, оказанных во время приема
func (s *Store) AddAppointmentServices(ctx context.Context, services []model.AppointmentService) error {
	if len(services) == 0 {
		return nil
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для добавления оказанных во время приема услуг: %w", err)
	}
	defer tx.Rollback()

	insert := s.builder.
		Insert("appointment_materials").
		Columns("visit_id", "material_id", "quantity_used")

	for _, m := range services {
		insert = insert.Values(m.VisitID, m.ServiceID, m.Quantity)
	}

	query, args, err := insert.ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать SQL-запрос для добавления оказанных во время приема услуг: %w", err)
	}

	_, err = tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос добавления оказанных во время приема услуг: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для добавления оказанных во время приема услуг: %w", err)
	}

	return nil
}
