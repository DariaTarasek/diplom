package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"github.com/DariaTarasek/diplom/services/admin/sharederrors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	minDurationMinutes = 10
	maxDurationMinutes = 180
)

func (s *AdminService) UpdateClinicSchedule(ctx context.Context, schedule []model.ClinicWeeklySchedule) error {
	var reqSchedule []*storagepb.WeeklyClinicSchedule
	for _, item := range schedule {
		if item.StartTime.Compare(item.EndTime) == 1 && item.IsDayOff {
			return fmt.Errorf("некорректный диапазон рабочего времени: %w", sharederrors.ErrInvalidValue)
		}
		if (item.SlotDurationMinutes < minDurationMinutes) || (item.SlotDurationMinutes > maxDurationMinutes) {
			return fmt.Errorf("некорректная продолжительность приема: %w", sharederrors.ErrInvalidValue)
		}
		if (item.EndTime.Sub(item.StartTime) < time.Duration(item.SlotDurationMinutes)*time.Minute) && item.IsDayOff {
			return fmt.Errorf("некорректная продолжительность приема: %w", sharederrors.ErrInvalidValue)
		}
		day := &storagepb.WeeklyClinicSchedule{
			Id:                  int32(item.ID),
			Weekday:             int32(item.Weekday),
			StartTime:           timestamppb.New(item.StartTime),
			EndTime:             timestamppb.New(item.EndTime),
			SlotDurationMinutes: int32(item.SlotDurationMinutes),
			IsDayOff:            !item.IsDayOff,
		}
		reqSchedule = append(reqSchedule, day)
	}
	_, err := s.StorageClient.Client.UpdateClinicWeeklySchedule(ctx, &storagepb.UpdateClinicWeeklyScheduleRequest{ClinicSchedule: reqSchedule})
	if err != nil {
		return fmt.Errorf("не удалось обновить постоянное расписание клиники: %w", err)
	}
	return nil
}

func (s *AdminService) UpdateDoctorSchedule(ctx context.Context, schedule []model.DoctorWeeklySchedule) error {
	var reqSchedule []*storagepb.WeeklyDoctorSchedule
	for _, item := range schedule {
		if item.StartTime.Compare(item.EndTime) == 1 && item.IsDayOff {
			return fmt.Errorf("некорректный диапазон рабочего времени: %w", sharederrors.ErrInvalidValue)
		}
		if (item.SlotDurationMinutes < minDurationMinutes) || (item.SlotDurationMinutes > maxDurationMinutes) {
			return fmt.Errorf("некорректная продолжительность приема: %w", sharederrors.ErrInvalidValue)
		}
		if (item.EndTime.Sub(item.StartTime) < time.Duration(item.SlotDurationMinutes)*time.Minute) && item.IsDayOff {
			return fmt.Errorf("некорректная продолжительность приема: %w", sharederrors.ErrInvalidValue)
		}
		day := &storagepb.WeeklyDoctorSchedule{
			Id:                  int32(item.ID),
			DoctorId:            int32(item.DoctorID),
			Weekday:             int32(item.Weekday),
			StartTime:           timestamppb.New(item.StartTime),
			EndTime:             timestamppb.New(item.EndTime),
			SlotDurationMinutes: int32(item.SlotDurationMinutes),
			IsDayOff:            !item.IsDayOff,
		}
		reqSchedule = append(reqSchedule, day)
	}
	_, err := s.StorageClient.Client.UpdateDoctorWeeklySchedule(ctx, &storagepb.UpdateDoctorWeeklyScheduleRequest{DoctorSchedule: reqSchedule})
	if err != nil {
		return fmt.Errorf("не удалось обновить постоянное расписание врача: %w", err)
	}
	return nil
}

func (s *AdminService) AddClinicDailyOverride(ctx context.Context, override model.ClinicDailyOverride) error {
	schedule, err := s.StorageClient.Client.GetClinicWeeklySchedule(ctx, &storagepb.EmptyRequest{})
	slot := schedule.ClinicSchedule[0].SlotDurationMinutes
	reqOverride := &storagepb.AddClinicDailyOverrideRequest{
		Date:                timestamppb.New(override.Date),
		StartTime:           timestamppb.New(override.StartTime),
		EndTime:             timestamppb.New(override.EndTime),
		SlotDurationMinutes: slot,
		IsDayOff:            override.IsDayOff,
	}

	_, err = s.StorageClient.Client.AddClinicDailyOverride(ctx, reqOverride)
	if err != nil {
		return fmt.Errorf("не удалось добавить переопределение дня работы клиники через gRPC: %w", err)
	}

	return nil
}

func (s *AdminService) AddDoctorDailyOverride(ctx context.Context, override model.DoctorDailyOverride) error {
	schedule, err := s.StorageClient.Client.GetDoctorWeeklySchedule(ctx, &storagepb.GetScheduleByDoctorIdRequest{DoctorId: int32(override.DoctorId)})
	slot := schedule.DoctorSchedule[0].SlotDurationMinutes
	reqOverride := &storagepb.AddDoctorDailyOverrideRequest{
		DoctorId:            int32(override.DoctorId),
		Date:                timestamppb.New(override.Date),
		StartTime:           timestamppb.New(override.StartTime),
		EndTime:             timestamppb.New(override.EndTime),
		SlotDurationMinutes: slot,
		IsDayOff:            override.IsDayOff,
	}

	_, err = s.StorageClient.Client.AddDoctorDailyOverride(ctx, reqOverride)
	if err != nil {
		return fmt.Errorf("не удалось добавить переопределение дня работы врача через gRPC: %w", err)
	}

	return nil
}
