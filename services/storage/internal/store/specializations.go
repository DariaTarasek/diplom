package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetSpecializations Получение списка всех специализаций
func (s *Store) GetSpecializations(ctx context.Context) ([]model.Specialization, error) {
	query, args, err := s.builder.
		Select("*").
		From("specializations").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка специализаций: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var specs []model.Specialization
	err = s.db.SelectContext(dbCtx, &specs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка специализаций: %w", err)
	}

	return specs, nil
}

// GetSpecializationByID Получение специализации по ID
func (s *Store) GetSpecializationByID(ctx context.Context, id model.SpecID) (model.Specialization, error) {
	query, args, err := s.builder.
		Select("*").
		From("specializations").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Specialization{}, fmt.Errorf("не удалось сформировать запрос для получения специализации: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var spec model.Specialization
	err = s.db.GetContext(dbCtx, &spec, query, args...)
	if err != nil {
		return model.Specialization{}, fmt.Errorf("не удалось выполнить запрос для получения специализации: %w", err)
	}

	return spec, nil
}

// AddSpecialization Добавление новой специализации
func (s *Store) AddSpecialization(ctx context.Context, spec model.Specialization) (model.SpecID, error) {
	fields := map[string]any{
		"name": spec.Name,
	}
	query, args, err := s.builder.
		Insert("specializations").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для добавления новой специализации: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var specID model.SpecID
	err = s.db.QueryRowxContext(dbCtx, query, args...).Scan(&specID)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для добавления новой специализации: %w", err)
	}
	return specID, nil
}

// UpdateSpecialization Изменение данных специализации
func (s *Store) UpdateSpecialization(ctx context.Context, id model.SpecID, spec model.Specialization) error {
	fields := map[string]any{
		"name": spec.Name,
	}
	query, args, err := s.builder.
		Update("specializations").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных специализации: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных специализации: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных специализации: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных специализации: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные специализации не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных специализации: %w", err)
	}
	return nil
}

// DeleteSpecialization Удаление специализации
func (s *Store) DeleteSpecialization(ctx context.Context, id model.SpecID) error {
	query, args, err := s.builder.
		Delete("specializations").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления специализации: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления специализации: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления специализации: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления специализации: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("специализация не удалена")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления специализации: %w", err)
	}
	return nil
}
