package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	_ "github.com/DariaTarasek/diplom/services/admin/sharederrors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sort"
	"strings"
	"time"
)

func NormalizeWord(input string) string {
	caser := cases.Title(language.Russian)
	return caser.String(strings.ToLower(strings.TrimSpace(input)))
}

func (s *AdminService) GetPatients(ctx context.Context) ([]model.Patient, error) {
	resp, err := s.StorageClient.Client.GetPatients(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список пациентов: %w", err)
	}
	var patients []model.Patient
	for _, item := range resp.Patients {
		patient := model.Patient{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Gender:      item.Gender,
			BirthDate:   item.BirthDate.AsTime().Format("02-01-2006"),
		}
		patients = append(patients, patient)
	}
	sort.Slice(patients, func(i, j int) bool {
		return patients[i].SecondName < patients[j].SecondName
	})

	return patients, nil
}

func (s *AdminService) UpdatePatient(ctx context.Context, patient model.Patient) error {
	bDate, err := time.Parse("2006-01-02", patient.BirthDate)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать дату рождения: %w", err)
	}
	_, err = s.StorageClient.Client.UpdatePatient(ctx, &storagepb.UpdatePatientRequest{
		UserId:      int32(patient.ID),
		FirstName:   NormalizeWord(patient.FirstName),
		Surname:     NormalizeWord(patient.Surname),
		SecondName:  NormalizeWord(patient.SecondName),
		Email:       strings.ToLower(patient.Email),
		BirthDate:   timestamppb.New(bDate),
		PhoneNumber: patient.PhoneNumber,
		Gender:      patient.Gender,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить данные пациента: %w", err)
	}

	return nil
}
