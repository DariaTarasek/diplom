package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"sort"
	"strings"
)

const (
	SuperAdminRole = 1
	AdminRole      = 2
)

func fetchRole(role int) string {
	switch role {
	case AdminRole:
		return "Администратор"
	case SuperAdminRole:
		return "Старший администратор"
	default:
		return "Администратор"
	}
}

func RoleToId(role string) int {
	switch role {
	case "Администратор":
		return AdminRole
	case "Старший администратор":
		return SuperAdminRole
	default:
		return AdminRole
	}
}

func (s *AdminService) GetAdmins(ctx context.Context) ([]model.Admin, error) {
	resp, err := s.StorageClient.Client.GetAdmins(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список админов: %w", err)
	}
	var admins []model.Admin
	for _, item := range resp.Admins {
		roleResp, err := s.StorageClient.Client.GetUserRole(ctx, &storagepb.GetUserRoleRequest{UserId: item.UserId})
		role := roleResp.Role
		if err != nil {
			return nil, fmt.Errorf("не удалось получить роль администратора: %w", err)
		}
		admin := model.Admin{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Gender:      item.Gender,
			Role:        fetchRole(int(role)),
		}
		admins = append(admins, admin)
	}
	sort.Slice(admins, func(i, j int) bool {
		return admins[i].SecondName < admins[j].SecondName
	})

	return admins, nil
}

func (s *AdminService) UpdateAdmin(ctx context.Context, admin model.Admin) error {
	_, err := s.StorageClient.Client.UpdateAdmin(ctx, &storagepb.UpdateAdminRequest{
		UserId:      int32(admin.ID),
		FirstName:   NormalizeWord(admin.FirstName),
		SecondName:  NormalizeWord(admin.SecondName),
		Surname:     NormalizeWord(admin.Surname),
		PhoneNumber: admin.PhoneNumber,
		Email:       strings.ToLower(admin.Email),
		Gender:      admin.Gender,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить данные администратора: %w", err)
	}
	_, err = s.StorageClient.Client.UpdateAdminRole(ctx, &storagepb.UpdateAdminRoleRequest{
		UserId: int32(admin.ID),
		RoleId: int32(RoleToId(admin.Role)),
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить роль администратора: %w", err)
	}
	return nil
}
