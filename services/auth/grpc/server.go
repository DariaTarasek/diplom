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
