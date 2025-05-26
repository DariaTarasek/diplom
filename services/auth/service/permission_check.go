package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	storagepb "github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

func (s *AuthService) PermissionCheck(ctx context.Context, token string, permissionID model.PermissionID) error {
	userID, err := utils.ParseToken(token)
	if err != nil {
		return fmt.Errorf("не удалось расшифровать токен: %w", err)
	}
	userRole, err := s.StorageClient.Client.GetUserRole(ctx, &storagepb.GetUserRoleRequest{UserId: userID})
	if err != nil {
		return fmt.Errorf("не удалось получить роль пользователя: %w", err)
	}
	_, err = s.StorageClient.Client.GetRolePermission(ctx, &storagepb.GetRolePermissionRequest{
		RoleId: userRole.Role,
		PermId: int32(permissionID),
	})
	if err != nil {
		return fmt.Errorf("не удалось получить право роли: %w", err)
	}
	return nil
}
