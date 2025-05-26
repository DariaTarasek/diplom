package grpcserver

import (
	"context"
	"github.com/DariaTarasek/diplom/services/auth/model"
	pb "github.com/DariaTarasek/diplom/services/auth/proto/auth"
	"github.com/DariaTarasek/diplom/services/auth/service"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	Service *service.AuthService
}

func (s *Server) EmployeeRegister(ctx context.Context, req *pb.EmployeeRegisterRequest) (*pb.EmployeeRegisterResponse, error) {
	modelUser := model.User{
		Login:    &req.User.Login,
		Password: &req.User.Password,
	}

	var id int
	var err error

	role := req.Employee.Role
	switch role {
	case model.DoctorRole:
		exp := int(req.Employee.Experience)

		modelDoctor := model.Doctor{
			FirstName:   req.Employee.FirstName,
			SecondName:  req.Employee.SecondName,
			Surname:     &req.Employee.Surname,
			PhoneNumber: &req.Employee.PhoneNumber,
			Email:       req.Employee.Email,
			Education:   &req.Employee.Education,
			Experience:  &exp,
			Gender:      req.Employee.Gender,
		}

		id, err = s.Service.DoctorRegister(ctx, modelUser, modelDoctor)
		if err != nil {
			return nil, err
		}
	case model.AdminRole:
		modelAdmin := model.Admin{
			FirstName:   req.Employee.FirstName,
			SecondName:  req.Employee.SecondName,
			Surname:     &req.Employee.Surname,
			PhoneNumber: &req.Employee.PhoneNumber,
			Email:       req.Employee.Email,
			Gender:      req.Employee.Gender,
		}

		id, err = s.Service.AdminRegister(ctx, modelUser, modelAdmin, false)
		if err != nil {
			return nil, err
		}
	case model.SuperAdminRole:
		modelAdmin := model.Admin{
			FirstName:   req.Employee.FirstName,
			SecondName:  req.Employee.SecondName,
			Surname:     &req.Employee.Surname,
			PhoneNumber: &req.Employee.PhoneNumber,
			Email:       req.Employee.Email,
			Gender:      req.Employee.Gender,
		}

		id, err = s.Service.AdminRegister(ctx, modelUser, modelAdmin, true)
		if err != nil {
			return nil, err
		}
	}

	return &pb.EmployeeRegisterResponse{
		UserId: int32(id),
	}, nil
}

func (s *Server) PatientRegister(ctx context.Context, req *pb.PatientRegisterRequest) (*pb.PatientRegisterResponse, error) {
	modelUser := model.User{
		Login:    &req.User.Login,
		Password: &req.User.Password,
	}

	patient := model.Patient{
		FirstName:   req.Patient.FirstName,
		SecondName:  req.Patient.SecondName,
		Surname:     &req.Patient.Surname,
		PhoneNumber: &req.Patient.PhoneNumber,
		Email:       &req.Patient.Email,
		BirthDate:   req.Patient.BirthDate.AsTime(),
		Gender:      req.Patient.Gender,
	}

	id, err := s.Service.PatientRegister(ctx, modelUser, patient)
	if err != nil {
		return nil, err
	}

	return &pb.PatientRegisterResponse{
		UserId: int32(id),
	}, nil
}

func (s *Server) PatientRegisterInClinic(ctx context.Context, req *pb.PatientRegisterInClinicRequest) (*pb.PatientRegisterInClinicResponse, error) {
	modelUser := model.User{
		Login:    &req.User.Login,
		Password: &req.User.Password,
	}

	modelPatient := model.Patient{
		FirstName:   req.Patient.FirstName,
		SecondName:  req.Patient.SecondName,
		Surname:     &req.Patient.Surname,
		PhoneNumber: &req.Patient.PhoneNumber,
		Email:       &req.Patient.Email,
		BirthDate:   req.Patient.BirthDate.AsTime(),
		Gender:      req.Patient.Gender,
	}

	id, err := s.Service.PatientRegisterInClinic(ctx, modelUser, modelPatient)
	if err != nil {
		return nil, err
	}

	return &pb.PatientRegisterInClinicResponse{
		UserId: int32(id),
	}, nil
}

func (s *Server) EmployeePasswordRecovery(ctx context.Context, req *pb.EmployeePasswordRecoveryRequest) (*pb.DefaultResponse, error) {
	err := s.Service.EmployeePasswordRecovery(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) PatientPasswordRecovery(ctx context.Context, req *pb.PatientPasswordRecoveryRequest) (*pb.DefaultResponse, error) {
	err := s.Service.PatientPasswordRecovery(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}
