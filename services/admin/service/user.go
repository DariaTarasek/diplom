package service

import (
	"context"
	"fmt"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"regexp"
	"strings"
)

func (s *AdminService) DeleteUser(ctx context.Context, id int) error {
	_, err := s.StorageClient.Client.DeleteUser(ctx, &storagepb.DeleteRequest{Id: int32(id)})
	if err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}
	return nil
}

func (s *AdminService) UpdateEmployeeLogin(ctx context.Context, id int, login string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(login) || len(login) > 256 {
		return fmt.Errorf("некорректный адрес электронной почты: %s", login)
	}

	_, err := s.StorageClient.Client.UpdateUserLogin(ctx, &storagepb.UpdateUserLoginRequest{
		UserId: int32(id),
		Login:  strings.ToLower(login),
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить логин сотрудника: %w", err)
	}

	return nil
}

func (s *AdminService) UpdatePatientLogin(ctx context.Context, id int, login string) error {
	phoneRegex := regexp.MustCompile(`^7\d{10}$`)
	if !phoneRegex.MatchString(login) {
		return fmt.Errorf("некорректный номер телефона: %s", login)
	}
	_, err := s.StorageClient.Client.UpdateUserLogin(ctx, &storagepb.UpdateUserLoginRequest{
		UserId: int32(id),
		Login:  login,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить логин пациента: %w", err)
	}

	return nil
}
