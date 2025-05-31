package grpcserver

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	pb "github.com/DariaTarasek/diplom/services/auth/proto/auth"
	"github.com/DariaTarasek/diplom/services/auth/service"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *Server) RequestCode(ctx context.Context, req *pb.GenerateCodeRequest) (*pb.DefaultResponse, error) {
	err := s.Service.RequestCode(ctx, req.Phone)
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить код подтверждения: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.DefaultResponse, error) {
	err := s.Service.VerifyCode(ctx, req.Phone, req.Code)
	if err != nil {
		return nil, fmt.Errorf("не удалось подтвердить код: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, role, err := s.Service.UserAuth(ctx, model.User{
		Login:    &req.Login,
		Password: &req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{
		Token: token,
		Role:  role,
	}, nil
}

func (s *Server) PermissionCheck(ctx context.Context, req *pb.PermissionCheckRequest) (*pb.DefaultResponse, error) {
	err := s.Service.PermissionCheck(ctx, req.Token, model.PermissionID(req.PermId))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetPatient(ctx context.Context, req *pb.GetPatientRequest) (*pb.GetPatientResponse, error) {
	patient, err := s.Service.GetPatientByID(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	pbPatient := &pb.PatientData{
		UserId:      int32(patient.ID),
		FirstName:   patient.FirstName,
		SecondName:  patient.SecondName,
		Surname:     deref(patient.Surname),
		Email:       deref(patient.Email),
		BirthDate:   timestamppb.New(patient.BirthDate),
		PhoneNumber: deref(patient.PhoneNumber),
		Gender:      patient.Gender,
	}
	return &pb.GetPatientResponse{Patient: pbPatient}, nil
}

func (s *Server) GetUserID(ctx context.Context, req *pb.GetUserIDRequest) (*pb.GetUserIDResponse, error) {
	userId, err := s.Service.GetUserID(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserIDResponse{UserId: int32(userId)}, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
