package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetRoles Получение списка всех ролей
func (s *Store) GetRoles(ctx context.Context) ([]model.Role, error) {
	query, args, err := s.builder.
		Select("*").
		From("roles").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка ролей: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var roles []model.Role
	err = s.db.SelectContext(dbCtx, &roles, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка ролей: %w", err)
	}

	return roles, nil
}

// GetRoleByID Получение ролей по ID
func (s *Store) GetRoleByID(ctx context.Context, id model.RoleID) (model.Role, error) {
	query, args, err := s.builder.
		Select("*").
		From("roles").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Role{}, fmt.Errorf("не удалось сформировать запрос для получения роли: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var role model.Role
	err = s.db.GetContext(dbCtx, &role, query, args...)
	if err != nil {
		return model.Role{}, fmt.Errorf("не удалось выполнить запрос для получения роли: %w", err)
	}

	return role, nil
}

// GetPermissions Получение списка всех прав
func (s *Store) GetPermissions(ctx context.Context) ([]model.Permission, error) {
	query, args, err := s.builder.
		Select("*").
		From("permissions").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка прав: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var permissions []model.Permission
	err = s.db.SelectContext(dbCtx, &permissions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка ролей: %w", err)
	}

	return permissions, nil
}

// GetPermissionByID Получение права по ID
func (s *Store) GetPermissionByID(ctx context.Context, id model.PermissionID) (model.Permission, error) {
	query, args, err := s.builder.
		Select("*").
		From("permissions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Permission{}, fmt.Errorf("не удалось сформировать запрос для получения права: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var permission model.Permission
	err = s.db.GetContext(dbCtx, &permission, query, args...)
	if err != nil {
		return model.Permission{}, fmt.Errorf("не удалось выполнить запрос для получения права: %w", err)
	}

	return permission, nil
}

// GetAllRolePermissions Получение списка всех прав роли
func (s *Store) GetAllRolePermissions(ctx context.Context, id model.RoleID) ([]model.RolePermissions, error) {
	query, args, err := s.builder.
		Select("*").
		From("role_permission").
		Where(squirrel.Eq{"role_id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка прав для роли: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var rolePermissions []model.RolePermissions
	err = s.db.SelectContext(dbCtx, &rolePermissions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка прав для роли: %w", err)
	}

	return rolePermissions, nil
}

// GetRolePermission Получение конкретного права роли
func (s *Store) GetRolePermission(ctx context.Context, roleID model.RoleID, permID model.PermissionID) (model.RolePermissions, error) {
	query, args, err := s.builder.
		Select("*").
		From("role_permission").
		Where(squirrel.Eq{"role_id": roleID, "permission_id": permID}).
		ToSql()
	if err != nil {
		return model.RolePermissions{}, fmt.Errorf("не удалось сформировать запрос для получения права роли: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var rolePermission model.RolePermissions
	err = s.db.GetContext(dbCtx, &rolePermission, query, args...)
	if err != nil {
		return model.RolePermissions{}, fmt.Errorf("не удалось выполнить запрос для получения права роли: %w", err)
	}

	return rolePermission, nil
}

// GetUsersByRole Получение списка пользователей с конкретной ролью
func (s *Store) GetUsersByRole(ctx context.Context, roleID model.RoleID) ([]model.UserRole, error) {
	query, args, err := s.builder.
		Select("*").
		From("user_role").
		Where(squirrel.Eq{"role_id": roleID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка пользователей с данной ролью: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var usersWithRole []model.UserRole
	err = s.db.SelectContext(dbCtx, &usersWithRole, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка пользователей с данной: %w", err)
	}

	return usersWithRole, nil
}

// GetRoleByUser Получение роли пользователя
func (s *Store) GetRoleByUser(ctx context.Context, userID model.UserID) (model.UserRole, error) {
	query, args, err := s.builder.
		Select("*").
		From("user_role").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return model.UserRole{}, fmt.Errorf("не удалось сформировать запрос для получения роли пользователя: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var userRole model.UserRole
	err = s.db.GetContext(dbCtx, &userRole, query, args...)
	if err != nil {
		return model.UserRole{}, fmt.Errorf("не удалось выполнить запрос для получения роли пользователя: %w", err)
	}

	return userRole, nil
}

// AddUserRole Добавление роли пользователю
func (s *Store) AddUserRole(ctx context.Context, userID model.UserID, roleID model.RoleID) error {
	fields := map[string]any{
		"user_id": userID,
		"role_id": roleID,
	}
	query, args, err := s.builder.
		Insert("user_role").
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления роли пользователю: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	_, err = s.db.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления роли пользователю: %w", err)
	}
	return nil
}

// UpdateUserRole Изменение роли пользователя
func (s *Store) UpdateUserRole(ctx context.Context, userID model.UserID, roleID model.RoleID) error {
	fields := map[string]any{
		"role_id": roleID,
	}
	query, args, err := s.builder.
		Update("user_role").
		Where(squirrel.Eq{"user_id": userID}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для обновления роли пользователя: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для изменения роли пользователя: %w", err)
	}

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для изменения роли пользователя: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после изменения роли пользователя: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("роль пользователя не изменена")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для изменения роли пользователя: %w", err)
	}
	return nil
}

// DeleteUserRole Удаление роли пользователя
func (s *Store) DeleteUserRole(ctx context.Context, userID model.UserID) error {
	query, args, err := s.builder.
		Delete("user_role").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления роли пользователя: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления роли пользователя: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления роли пользователя: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество измененных строк после удаления роли пользователя: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("роль пользователя не удалена")
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления роли пользователя: %w", err)
	}
	return nil
}
