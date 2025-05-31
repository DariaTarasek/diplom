package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	storagepb "github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/sharederrors"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

func (s *AuthService) UserAuth(ctx context.Context, user model.User) (string, string, error) {
	storageUser, err := s.StorageClient.Client.GetUserByLogin(ctx, &storagepb.GetUserByLoginRequest{Login: deref(user.Login)})
	if err != nil {
		return "", "", fmt.Errorf("не удалось получить пользователя из базы: %w", err)
	}
	fmt.Println(storageUser)
	password := deref(user.Password)
	fmt.Println(password)
	if err != nil {
		return "", "", fmt.Errorf("не удалось хешировать полученный пароль: %w", err)
	}
	if err := utils.ComparePasswords(storageUser.Password, password); err != nil {
		return "", "", fmt.Errorf("введен неверный пароль: %w", sharederrors.ErrPasswordInvalid)
	}

	token, err := utils.GenerateToken(int(storageUser.Id))
	if err != nil {
		return "", "", fmt.Errorf("не удалось сгенерировать токен: %w", err)
	}
	role, err := s.StorageClient.Client.GetUserRole(ctx, &storagepb.GetUserRoleRequest{UserId: storageUser.Id})
	if err != nil {
		return "", "", fmt.Errorf("не удалось получить роль пользователя: %w", err)
	}

	return token, fetchRole(int(role.Role)), nil

}

func (s *AuthService) GetUserID(ctx context.Context, token string) (model.UserID, error) {
	userID, err := utils.ParseToken(token)
	if err != nil {
		return 0, fmt.Errorf("не удалось разобрать токен: %w", err)
	}
	return model.UserID(userID), nil
}

func fetchRole(role int) string {
	switch role {
	case model.DoctorRole:
		return "doctor"
	case model.AdminRole:
		return "admin"
	case model.SuperAdminRole:
		return "superadmin"
	default:
		return "patient"
	}
}
