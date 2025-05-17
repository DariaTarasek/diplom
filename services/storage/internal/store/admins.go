package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetAdmins Получение списка всех администраторов
func (s *Store) GetAdmins(ctx context.Context) ([]model.Admin, error) {
	query, args, err := s.builder.
		Select("*").
		From("administrators").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка администраторов: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var admins []model.Admin
	err = s.db.SelectContext(dbCtx, &admins, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка администраторов: %w", err)
	}

	return admins, nil
}

// GetAdminByID Получение администратора по id
func (s *Store) GetAdminByID(ctx context.Context, id model.UserID) (model.Admin, error) {
	query, args, err := s.builder.
		Select("*").
		From("administrators").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return model.Admin{}, fmt.Errorf("не удалось сформировать запрос для получения администратора по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var admin model.Admin
	err = s.db.GetContext(dbCtx, &admin, query, args...)
	if err != nil {
		return model.Admin{}, fmt.Errorf("не удалось выполнить запрос для получения администратора по id: %w", err)
	}

	return admin, nil
}

// AddAdmin Добавление нового администратора
func (s *Store) AddAdmin(ctx context.Context, admin model.Admin) error {
	fields := map[string]any{
		"user_id":      admin.ID,
		"first_name":   admin.FirstName,
		"second_name":  admin.SecondName,
		"surname":      admin.Surname,
		"phone_number": admin.PhoneNumber,
		"email":        admin.Email,
		"gender":       admin.Gender,
	}
	query, args, err := s.builder.
		Insert("administrators").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления нового администратора: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления нового администратора: %w", err)
	}
	return nil
}

// UpdateAdmin Изменение данных администратора
func (s *Store) UpdateAdmin(ctx context.Context, id model.UserID, admin model.Admin) error {
	fields := map[string]any{
		"first_name":   admin.FirstName,
		"second_name":  admin.SecondName,
		"surname":      admin.Surname,
		"phone_number": admin.PhoneNumber,
		"email":        admin.Email,
		"gender":       admin.Gender,
	}
	query, args, err := s.builder.
		Update("administrators").
		Where(squirrel.Eq{"user_id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных администратора: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных администратора: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных администратора: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных администратора: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные администратора не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных администратора: %w", err)
	}
	return nil
}

// DeleteAdmin Удаление администратора по id
func (s *Store) DeleteAdmin(ctx context.Context, id model.UserID) error {
	query, args, err := s.builder.
		Delete("administrators").
		Where(squirrel.Eq{"user_id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления администратора: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления администратора: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления администратора: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления администратора: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("администратор не удален")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления администратора: %w", err)
	}
	return nil
}
