package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/clients"
	"github.com/DariaTarasek/diplom/services/auth/model"
	"github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

type AuthService struct {
	StorageClient *clients.StorageClient
}

func NewAuthService(client *clients.StorageClient) *AuthService {
	return &AuthService{StorageClient: client}
}

func (s *AuthService) DoctorRegister(ctx context.Context, user model.User, doctor model.Doctor) (int, error) {
	var hashedPassword string
	var err error
	if user.Password != nil {
		hashedPassword, err = utils.HashPassword(*user.Password)
		if err != nil {
			return 0, fmt.Errorf("не удалось захешировать пароль: %w", err)
		}
	}

	reqUser := &storagepb.AddUserRequest{
		Login:    *user.Login,
		Password: hashedPassword,
	}

	respUser, err := s.StorageClient.Client.AddUser(ctx, reqUser)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить пользователя через gRPC: %w", err)
	}

	reqDoctor := &storagepb.AddDoctorRequest{
		UserId:      respUser.UserId,
		FirstName:   doctor.FirstName,
		SecondName:  doctor.SecondName,
		Surname:     *doctor.Surname,
		PhoneNumber: *doctor.PhoneNumber,
		Email:       doctor.Email,
		Education:   *doctor.Education,
		Experience:  int32(*doctor.Experience),
		Gender:      doctor.Gender,
	}

	_, err = s.StorageClient.Client.AddDoctor(ctx, reqDoctor)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить врача через gRPC: %w", err)
	}

	return int(respUser.UserId), nil
}
