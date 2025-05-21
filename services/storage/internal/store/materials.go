package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetMaterials Получение списка всех материалов
func (s *Store) GetMaterials(ctx context.Context) ([]model.Material, error) {
	query, args, err := s.builder.
		Select("*").
		From("materials").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка материалов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var materials []model.Material
	err = s.db.SelectContext(dbCtx, &materials, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка материалов: %w", err)
	}

	return materials, nil
}

// GetMaterialByID Получение материалa по id
func (s *Store) GetMaterialByID(ctx context.Context, id model.MaterialID) (model.Material, error) {
	query, args, err := s.builder.
		Select("*").
		From("materials").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Material{}, fmt.Errorf("не удалось сформировать запрос для получения материалa по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var material model.Material
	err = s.db.GetContext(dbCtx, &material, query, args...)
	if err != nil {
		return model.Material{}, fmt.Errorf("не удалось выполнить запрос для получения материалa по id: %w", err)
	}

	return material, nil
}

// AddMaterial Добавление нового материалa
func (s *Store) AddMaterial(ctx context.Context, material model.Material) (model.MaterialID, error) {
	fields := map[string]any{
		"name":  material.Name,
		"price": material.Price,
	}
	query, args, err := s.builder.
		Insert("materials").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return model.MaterialID(0), fmt.Errorf("не удалось сформировать запрос для добавления нового материала: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var materialID model.MaterialID
	err = s.db.QueryRowxContext(dbCtx, query, args...).Scan(&materialID)
	if err != nil {
		return model.MaterialID(0), fmt.Errorf("не удалось выполнить запрос для добавления нового материала: %w", err)
	}

	return materialID, nil
}

// UpdateMaterial Изменение данных материала
func (s *Store) UpdateMaterial(ctx context.Context, id model.MaterialID, material model.Material) error {
	fields := map[string]any{
		"name":  material.Name,
		"price": material.Price,
	}
	query, args, err := s.builder.
		Update("materials").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных материала: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных материала: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных материала: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных материала: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные материала не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных материала: %w", err)
	}
	return nil
}

// DeleteMaterial Удаление материала по id
func (s *Store) DeleteMaterial(ctx context.Context, id model.MaterialID) error {
	query, args, err := s.builder.
		Delete("materials").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления материала: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления материала: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления материала: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления материала: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("материал не удален")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления материала: %w", err)
	}
	return nil
}
