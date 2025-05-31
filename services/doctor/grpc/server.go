package grpcserver

import (
	"context"
	pb "github.com/DariaTarasek/diplom/services/doctor/proto/doctor"
	"github.com/DariaTarasek/diplom/services/doctor/service"
)

type Server struct {
	pb.UnimplementedDoctorServiceServer
	Service *service.DoctorService
}

func (s *Server) GetTodayAppointments(ctx context.Context, request *pb.GetTodayAppointmentsRequest) (*pb.GetTodayAppointmentsResponse, error) {
	apps, err := s.Service.GetTodayAppointments(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	todays := make([]*pb.TodayAppointment, 0, len(apps))
	for _, app := range apps {
		today := &pb.TodayAppointment{
			Id:        int32(app.ID),
			Date:      app.Date,
			Time:      app.Time,
			PatientID: int32(app.PatientID),
			Patient:   app.Patient,
		}
		todays = append(todays, today)
	}
	return &pb.GetTodayAppointmentsResponse{Appointments: todays}, nil
}

func (s *Server) GetUpcomingAppointments(ctx context.Context, request *pb.GetUpcomingAppointmentsRequest) (*pb.GetUpcomingAppointmentsResponse, error) {
	apps, err := s.Service.GetUpcomingAppointments(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	var cells []*pb.ScheduleCell
	for date, row := range apps.Table {
		for timeStr, appointment := range row {
			cell := &pb.ScheduleCell{
				Date: date,
				Time: timeStr,
			}
			if appointment != nil {
				cell.Appointment = &pb.UpcomingAppointment{
					Id:        int64(appointment.ID),
					PatientId: int64(appointment.PatientID),
					Patient:   appointment.Patient,
				}
			}
			cells = append(cells, cell)
		}
	}

	resp := &pb.ScheduleTable{
		Dates: apps.Dates,
		Times: apps.Times,
		Table: cells,
	}

	return &pb.GetUpcomingAppointmentsResponse{Schedule: resp}, nil

}
