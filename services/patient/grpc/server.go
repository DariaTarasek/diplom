package grpcserver

import (
	"context"
	"github.com/DariaTarasek/diplom/services/patient/model"
	pb "github.com/DariaTarasek/diplom/services/patient/proto/patient"
	"github.com/DariaTarasek/diplom/services/patient/service"
)

type Server struct {
	pb.UnimplementedPatientServiceServer
	Service *service.PatientService
}

func (s *Server) GetAppointmentSlots(ctx context.Context, request *pb.GetAppointmentSlotsRequest) (*pb.GetAppointmentSlotsResponse, error) {
	slots, err := s.Service.MakeDoctorAppointmentSlots(ctx, int(request.DoctorId))
	if err != nil {
		return nil, err
	}
	var res []*pb.DaySlots
	for _, slot := range slots {
		res = append(res, &pb.DaySlots{
			Date:  slot.Label,
			Slots: slot.Slots,
		})
	}
	return &pb.GetAppointmentSlotsResponse{Slots: res}, nil
}

func (s *Server) AddAppointment(ctx context.Context, request *pb.AddAppointmentRequest) (*pb.DefaultResponse, error) {
	patientID := model.UserID(request.Appointment.PatientId)
	appointment := model.Appointment{
		DoctorID:           model.UserID(request.Appointment.DoctorId),
		PatientID:          &patientID,
		Date:               request.Appointment.Date.AsTime(),
		Time:               request.Appointment.Time.AsTime(),
		PatientSecondName:  request.Appointment.SecondName,
		PatientFirstName:   request.Appointment.FirstName,
		PatientSurname:     &request.Appointment.Surname,
		PatientBirthDate:   request.Appointment.BirthDate.AsTime(),
		PatientGender:      request.Appointment.Gender,
		PatientPhoneNumber: request.Appointment.PhoneNumber,
	}
	err := s.Service.AddAppointment(ctx, appointment)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}
