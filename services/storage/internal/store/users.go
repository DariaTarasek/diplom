package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetUsers Получение списка всех пользователей
func (s *Store) GetUsers(ctx context.Context) ([]model.User, error) {
	query, args, err := s.builder.
		Select("*").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка пользователей: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var users []model.User
	err = s.db.SelectContext(dbCtx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка пользователей: %w", err)
	}

	return users, nil
}

// GetUserByID Получение пользователя по id
func (s *Store) GetUserByID(ctx context.Context, id model.UserID) (model.User, error) {
	query, args, err := s.builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("не удалось сформировать запрос для получения пользователя по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var user model.User
	err = s.db.GetContext(dbCtx, &user, query, args...)
	if err != nil {
		return model.User{}, fmt.Errorf("не удалось выполнить запрос для получения пользователя по id: %w", err)
	}

	return user, nil
}

// GetUserByLogin Получение пользователя по логину
func (s *Store) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	query, args, err := s.builder.
		Select("*").
		From("users").
		Where(squirrel.Eq{"login": login}).
		ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("не удалось сформировать запрос для получения пользователя по логину: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var user model.User
	err = s.db.GetContext(dbCtx, &user, query, args...)
	if err != nil {
		return model.User{}, fmt.Errorf("не удалось выполнить запрос для получения пользователя по логину: %w", err)
	}

	return user, nil
}

// AddUser Добавление нового пользователя
func (s *Store) AddUser(ctx context.Context, user model.User) (model.UserID, error) {
	fields := map[string]any{
		"login":    user.Login,
		"password": user.Password,
	}
	query, args, err := s.builder.
		Insert("users").
		SetMap(fields).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return model.UserID(0), fmt.Errorf("не удалось сформировать запрос для добавления нового пользователя: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var userID model.UserID
	err = s.db.QueryRowxContext(dbCtx, query, args...).Scan(&userID)
	if err != nil {
		return model.UserID(0), fmt.Errorf("не удалось выполнить запрос для добавления нового пользователя: %w", err)
	}

	return userID, nil
}

// UpdateUser Изменение данных пользователя
func (s *Store) UpdateUser(ctx context.Context, id model.UserID, user model.User) error {
	fields := map[string]any{
		"login":    user.Login,
		"password": user.Password,
	}
	query, args, err := s.builder.
		Update("users").
		Where(squirrel.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления данных пользователя: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения данных пользователя: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения данных пользователя: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения данных пользователя: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("данные пользователя не изменены")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения данных пользователя: %w", err)
	}
	return nil
}

// DeleteUser Удаление администратора по id
func (s *Store) DeleteUser(ctx context.Context, id model.UserID) error {
	query, args, err := s.builder.
		Delete("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления пользователя: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления пользователя: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления пользователя: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления пользователя: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь не удален")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления пользователя: %w", err)
	}
	return nil
}
