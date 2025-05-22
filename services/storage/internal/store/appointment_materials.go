package store

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetAppointmentMaterials Получение списка всех потраченных во время приема материалов
func (s *Store) GetAppointmentMaterials(ctx context.Context, id model.VisitID) ([]model.AppointmentMaterial, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_materials").
		Where(squirrel.Eq{"visit_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения потраченных во время приема материалов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var appMaterials []model.AppointmentMaterial
	err = s.db.SelectContext(dbCtx, &appMaterials, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения потраченных во время приема материалов: %w", err)
	}

	return appMaterials, nil
}

// AddAppointmentMaterials Добавление записей о материалах, использованных во время приема
func (s *Store) AddAppointmentMaterials(ctx context.Context, materials []model.AppointmentMaterial) error {
	if len(materials) == 0 {
		return nil
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, nil)
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для добавления материалов, использованных во время приема: %w", err)
	}
	defer tx.Rollback()

	insert := s.builder.
		Insert("appointment_materials").
		Columns("visit_id", "material_id", "quantity_used")

	for _, m := range materials {
		insert = insert.Values(m.VisitID, m.MaterialID, m.QuantityUsed)
	}

	query, args, err := insert.ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать SQL-запрос для добавления материалов, использованных во время приема: %w", err)
	}

	_, err = tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос добавления материалов, использованных во время приема: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для добавления материалов, использованных во время приема: %w", err)
	}

	return nil
}
