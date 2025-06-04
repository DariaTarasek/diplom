package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	storagepb "github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

func (s *AuthService) GetAdminByID(ctx context.Context, token string) (model.AdminWithRole, error) {
	userID, err := utils.ParseToken(token)
	if err != nil {
		return model.AdminWithRole{}, fmt.Errorf("не удалось получить id пользователя: %w", err)
	}
	admin, err := s.StorageClient.Client.GetAdminByID(ctx, &storagepb.GetByIDRequest{Id: userID})
	if err != nil {
		return model.AdminWithRole{}, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	userRole, err := s.StorageClient.Client.GetUserRole(ctx, &storagepb.GetUserRoleRequest{UserId: userID})
	if err != nil {
		return model.AdminWithRole{}, fmt.Errorf("не удалось получить роль администратора: %w", err)
	}
	return model.AdminWithRole{
		ID:          model.UserID(admin.Admin.UserId),
		FirstName:   admin.Admin.FirstName,
		SecondName:  admin.Admin.SecondName,
		Surname:     &admin.Admin.Surname,
		PhoneNumber: &admin.Admin.PhoneNumber,
		Email:       admin.Admin.Email,
		Gender:      admin.Admin.Gender,
		Role:        fetchRole(int(userRole.Role)),
	}, nil
}
