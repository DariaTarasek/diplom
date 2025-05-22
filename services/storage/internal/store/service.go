package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetServiceTypes Получение списка всех типов услуг
func (s *Store) GetServiceTypes(ctx context.Context) ([]model.ServiceType, error) {
	query, args, err := s.builder.
		Select("*").
		From("service_types").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка типов услуг: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var types []model.ServiceType
	err = s.db.SelectContext(dbCtx, &types, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка типов услуг: %w", err)
	}

	return types, nil
}

// GetServiceTypeByID Получение типа услуги по ID
func (s *Store) GetServiceTypeByID(ctx context.Context, id model.ServiceTypeID) (model.ServiceType, error) {
	query, args, err := s.builder.
		Select("*").
		From("service_types").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.ServiceType{}, fmt.Errorf("не удалось сформировать запрос для получения типа услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var serviceType model.ServiceType
	err = s.db.GetContext(dbCtx, &serviceType, query, args...)
	if err != nil {
		return model.ServiceType{}, fmt.Errorf("не удалось выполнить запрос для получения типа услуги: %w", err)
	}

	return serviceType, nil
}

// AddServiceType Добавление нового типа услуги
func (s *Store) AddServiceType(ctx context.Context, serviceType model.ServiceType) (model.ServiceTypeID, error) {
	fields := map[string]any{
		"name": serviceType.Name,
	}
	query, args, err := s.builder.
		Insert("service_types").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для добавления нового типа услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var id model.ServiceTypeID
	err = s.db.QueryRowContext(dbCtx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для добавления нового типа услуги: %w", err)
	}

	return id, nil
}

// UpdateServiceType Изменение типа услуги
func (s *Store) UpdateServiceType(ctx context.Context, id model.ServiceTypeID, serviceType model.ServiceType) error {
	fields := map[string]any{
		"name": serviceType.Name,
	}
	query, args, err := s.builder.
		Update("service_types").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для изменения типа услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения типа услуги: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения типа услуги: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения типа услуги: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("тип услуги не изменен")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения типа услуги: %w", err)
	}
	return nil
}

// DeleteServiceType Удаление типа услуги по id
func (s *Store) DeleteServiceType(ctx context.Context, id model.ServiceTypeID) error {
	query, args, err := s.builder.
		Delete("service_types").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления типа услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления типа услуги: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления типа услуги: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления типа услуги: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("тип услуги не удален")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления типа услуги: %w", err)
	}
	return nil
}

// GetServices Получение списка всех услуг
func (s *Store) GetServices(ctx context.Context) ([]model.Service, error) {
	query, args, err := s.builder.
		Select("*").
		From("services").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка услуг: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var services []model.Service
	err = s.db.SelectContext(dbCtx, &services, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка услуг: %w", err)
	}

	return services, nil
}

// GetServiceByID Получение услуги по ID
func (s *Store) GetServiceByID(ctx context.Context, id model.ServiceID) (model.Service, error) {
	query, args, err := s.builder.
		Select("*").
		From("services").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Service{}, fmt.Errorf("не удалось сформировать запрос для получения услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var service model.Service
	err = s.db.GetContext(dbCtx, &service, query, args...)
	if err != nil {
		return model.Service{}, fmt.Errorf("не удалось выполнить запрос для получения услуги: %w", err)
	}

	return service, nil
}

// AddService Добавление новой услуги
func (s *Store) AddService(ctx context.Context, service model.Service) (model.ServiceID, error) {
	fields := map[string]any{
		"name":  service.Name,
		"price": service.Price,
		"type":  service.Category,
	}
	query, args, err := s.builder.
		Insert("services").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("не удалось сформировать запрос для добавления новой услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var id model.ServiceID
	err = s.db.QueryRowContext(dbCtx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось выполнить запрос для добавления новой услуги: %w", err)
	}

	return id, nil
}

// UpdateService Изменение услуги
func (s *Store) UpdateService(ctx context.Context, id model.ServiceID, service model.Service) error {
	fields := map[string]any{
		"name":  service.Name,
		"price": service.Price,
		"type":  service.Category,
	}
	query, args, err := s.builder.
		Update("services").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для изменения услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения услуги: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения услуги: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения услуги: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("услуга не изменена")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения услуги: %w", err)
	}
	return nil
}

// DeleteService Удаление услуги по id
func (s *Store) DeleteService(ctx context.Context, id model.ServiceID) error {
	query, args, err := s.builder.
		Delete("services").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления услуги: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления услуги: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления услуги: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления услуги: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("услуга не удалена")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления услуги: %w", err)
	}
	return nil
}
