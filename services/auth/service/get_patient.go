package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	storagepb "github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

func (s *AuthService) GetPatientByID(ctx context.Context, token string) (model.Patient, error) {
	userID, err := utils.ParseToken(token)
	if err != nil {
		return model.Patient{}, fmt.Errorf("не удалось получить id пользователя: %w", err)
	}
	patient, err := s.StorageClient.Client.GetPatientByID(ctx, &storagepb.GetByIDRequest{Id: userID})
	if err != nil {
		return model.Patient{}, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	return model.Patient{
		ID:          model.UserID(patient.Patient.UserId),
		FirstName:   patient.Patient.FirstName,
		SecondName:  patient.Patient.SecondName,
		Surname:     &patient.Patient.Surname,
		Email:       &patient.Patient.Email,
		BirthDate:   patient.Patient.BirthDate.AsTime(),
		PhoneNumber: &patient.Patient.PhoneNumber,
		Gender:      patient.Patient.Gender,
	}, nil
}
