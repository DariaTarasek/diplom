package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

func (s *AuthService) EmployeePasswordRecovery(ctx context.Context, login string) error {
	resp, err := s.StorageClient.Client.GetUserByLogin(ctx, &storagepb.GetUserByLoginRequest{Login: login})
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	password := utils.GeneratePassword()
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("не удалось захешировать пароль: %w", err)
	}

	_, err = s.StorageClient.Client.UpdateUserPassword(ctx, &storagepb.UpdateUserPasswordRequest{
		Id:       resp.Id,
		Login:    resp.Login,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить пароль: %w", err)
	}

	message := fmt.Sprintf("Subject: Восстановление пароля\r\n\r\nВы запросили восстановление пароля в системе клиники.\nВаш новый пароль: %s", password)
	err = utils.SendPassword(resp.Login, password, message)
	if err != nil {
		return fmt.Errorf("не удалось отправить пароль на email: %w", err)
	}

	return nil
}

func (s *AuthService) PatientPasswordRecovery(ctx context.Context, login string) error {
	resp, err := s.StorageClient.Client.GetUserByLogin(ctx, &storagepb.GetUserByLoginRequest{Login: login})
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	password := utils.GeneratePassword()
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("не удалось захешировать пароль: %w", err)
	}

	_, err = s.StorageClient.Client.UpdateUserPassword(ctx, &storagepb.UpdateUserPasswordRequest{
		Id:       resp.Id,
		Login:    resp.Login,
		Password: hashedPassword,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить пароль: %w", err)
	}

	// отправка нового пароля по sms

	return nil
}
