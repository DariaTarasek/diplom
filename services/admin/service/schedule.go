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
