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

func (s *Server) UpdateDoctorWeeklySchedule(ctx context.Context, req *pb.UpdateDoctorWeeklyScheduleRequest) (*pb.DefaultResponse, error) {
	var reqSchedule []model.DoctorWeeklySchedule
	for _, item := range req.DoctorSchedule {
		day := model.DoctorWeeklySchedule{
			ID:                  int(item.Id),
			DoctorID:            int(item.DoctorId),
			Weekday:             int(item.Weekday),
			StartTime:           item.StartTime.AsTime(),
			EndTime:             item.EndTime.AsTime(),
			SlotDurationMinutes: int(item.SlotDurationMinutes),
			IsDayOff:            item.IsDayOff,
		}
		reqSchedule = append(reqSchedule, day)
	}
	err := s.Service.UpdateDoctorSchedule(ctx, reqSchedule)
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить расписание врача: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddClinicDailyOverride(ctx context.Context, req *pb.AddClinicDailyOverrideRequest) (*pb.DefaultResponse, error) {
	override := model.ClinicDailyOverride{
		Date:                req.Date.AsTime(),
		StartTime:           req.StartTime.AsTime(),
		EndTime:             req.EndTime.AsTime(),
		SlotDurationMinutes: int(req.SlotDurationMinutes),
		IsDayOff:            req.IsDayOff,
	}

	err := s.Service.AddClinicDailyOverride(ctx, override)
	if err != nil {
		return nil, err
	}

	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddDoctorDailyOverride(ctx context.Context, req *pb.AddDoctorDailyOverrideRequest) (*pb.DefaultResponse, error) {
	override := model.DoctorDailyOverride{
		DoctorId:            int(req.DoctorId),
		Date:                req.Date.AsTime(),
		StartTime:           req.StartTime.AsTime(),
		EndTime:             req.EndTime.AsTime(),
		SlotDurationMinutes: int(req.SlotDurationMinutes),
		IsDayOff:            req.IsDayOff,
	}

	err := s.Service.AddDoctorDailyOverride(ctx, override)
	if err != nil {
		return nil, err
	}

	return &pb.DefaultResponse{}, nil
}
