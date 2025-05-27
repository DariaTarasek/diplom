package grpc

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	pb "github.com/DariaTarasek/diplom/services/admin/proto/admin"
	"github.com/DariaTarasek/diplom/services/admin/service"
)

type Server struct {
	pb.UnimplementedAdminServiceServer
	Service *service.AdminService
}

func (s *Server) UpdateClinicWeeklySchedule(ctx context.Context, req *pb.UpdateClinicWeeklyScheduleRequest) (*pb.DefaultResponse, error) {
	var reqSchedule []model.ClinicWeeklySchedule
	for _, item := range req.ClinicSchedule {
		day := model.ClinicWeeklySchedule{
			ID:                  int(item.Id),
			Weekday:             int(item.Weekday),
			StartTime:           item.StartTime.AsTime(),
			EndTime:             item.EndTime.AsTime(),
			SlotDurationMinutes: int(item.SlotDurationMinutes),
			IsDayOff:            item.IsDayOff,
		}
		reqSchedule = append(reqSchedule, day)
	}
	err := s.Service.UpdateClinicSchedule(ctx, reqSchedule)
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить расписание клиники: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}
