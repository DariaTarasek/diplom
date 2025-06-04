package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	storagepb "github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
)

func (s *AuthService) GetDoctorByID(ctx context.Context, token string) (model.Doctor, error) {
	userID, err := utils.ParseToken(token)
	if err != nil {
		return model.Doctor{}, fmt.Errorf("не удалось получить id пользователя: %w", err)
	}
	doc, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: userID})
	if err != nil {
		return model.Doctor{}, fmt.Errorf("не удалось получить пользователя: %w", err)
	}
	specsResp, err := s.StorageClient.Client.GetSpecsByDoctorID(ctx, &storagepb.GetByIDRequest{Id: userID})
	specs := make([]int, 0, len(specsResp.SpecId))
	for _, spec := range specsResp.SpecId {
		specs = append(specs, int(spec))
	}
	exp := int(doc.Doctor.Experience)
	return model.Doctor{
		ID:          model.UserID(doc.Doctor.UserId),
		FirstName:   doc.Doctor.FirstName,
		SecondName:  doc.Doctor.SecondName,
		Surname:     &doc.Doctor.Surname,
		PhoneNumber: &doc.Doctor.PhoneNumber,
		Email:       doc.Doctor.Email,
		Education:   &doc.Doctor.Education,
		Experience:  &exp,
		Gender:      doc.Doctor.Gender,
		Specs:       specs,
	}, nil
}
