package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"slices"
	"sort"
)

func (s *AdminService) GetSpecs(ctx context.Context) ([]model.Spec, error) {
	resp, err := s.StorageClient.Client.GetAllSpecs(ctx, &storagepb.EmptyRequest{})
	var specs []model.Spec
	for _, item := range resp.Specs {
		spec := model.Spec{
			ID:   int(item.Id),
			Name: item.Name,
		}
		specs = append(specs, spec)
	}
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список специализаций врачей: %w", err)
	}
	return specs, nil
}

func (s *AdminService) GetDoctors(ctx context.Context) ([]model.Doctor, error) {
	resp, err := s.StorageClient.Client.GetDoctors(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список врачей: %w", err)
	}
	var doctors []model.Doctor
	for _, item := range resp.Doctors {
		specsResp, err := s.StorageClient.Client.GetDoctorSpecsByDoctorId(ctx, &storagepb.GetByIdRequest{Id: item.UserId})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить специальности врача: %w", err)
		}
		var specs []int
		for _, specItem := range specsResp.Specs {
			intSpecItem := int(specItem)
			specs = append(specs, intSpecItem)
		}
		doctor := model.Doctor{
			ID:          int(item.UserId),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Education:   item.Education,
			Experience:  int(item.Experience),
			Gender:      item.Gender,
			Specs:       specs,
		}
		doctors = append(doctors, doctor)
	}

	sort.Slice(doctors, func(i, j int) bool {
		return doctors[i].SecondName < doctors[j].SecondName
	})

	return doctors, nil
}

func (s *AdminService) UpdateDoctor(ctx context.Context, doctor model.Doctor) error {
	_, err := s.StorageClient.Client.UpdateDoctor(ctx, &storagepb.UpdateDoctorRequest{
		UserId:      int32(doctor.ID),
		FirstName:   doctor.FirstName,
		SecondName:  doctor.SecondName,
		Surname:     doctor.Surname,
		PhoneNumber: doctor.PhoneNumber,
		Email:       doctor.Email,
		Education:   doctor.Education,
		Experience:  int32(doctor.Experience),
		Gender:      doctor.Gender,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить данные врача: %w", err)
	}
	oldSpecsResp, err := s.StorageClient.Client.GetDoctorSpecsByDoctorId(ctx, &storagepb.GetByIdRequest{Id: int32(doctor.ID)})
	if err != nil {
		return fmt.Errorf("не удалось получить специализации врача: %w", err)
	}
	var oldSpecs []int
	for _, item := range oldSpecsResp.Specs {
		oldSpec := int(item)
		oldSpecs = append(oldSpecs, oldSpec)
	}
	for _, item := range oldSpecs {
		if !slices.Contains(doctor.Specs, item) {
			_, err = s.StorageClient.Client.DeleteDoctorSpec(ctx, &storagepb.DeleteDoctorSpecRequest{
				DoctorId: int32(doctor.ID),
				SpecId:   int32(item),
			})
			if err != nil {
				return fmt.Errorf("не удалось удалить старую специализацию врача: %w", err)
			}
		}
	}
	for _, item := range doctor.Specs {
		if !slices.Contains(oldSpecs, item) {
			_, err = s.StorageClient.Client.AddDoctorSpec(ctx, &storagepb.AddDoctorSpecRequest{
				DoctorId: int32(doctor.ID),
				SpecId:   int32(item),
			})
			if err != nil {
				return fmt.Errorf("не удалось добавить новую специализацию врача: %w", err)
			}
		}
	}
	return nil
}
