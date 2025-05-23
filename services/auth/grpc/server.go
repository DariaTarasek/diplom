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

func (s *Server) DoctorRegister(ctx context.Context, req *pb.DoctorRegisterRequest) (*pb.DoctorRegisterResponse, error) {
	modelUser := model.User{
		Login:    &req.User.Login,
		Password: &req.User.Password,
	}

	exp := int(req.Doctor.Experience)

	modelDoctor := model.Doctor{
		FirstName:   req.Doctor.FirstName,
		SecondName:  req.Doctor.SecondName,
		Surname:     &req.Doctor.Surname,
		PhoneNumber: &req.Doctor.PhoneNumber,
		Email:       req.Doctor.Email,
		Education:   &req.Doctor.Education,
		Experience:  &exp,
		Gender:      req.Doctor.Gender,
	}

	id, err := s.Service.DoctorRegister(ctx, modelUser, modelDoctor)
	if err != nil {
		return nil, err
	}

	return &pb.DoctorRegisterResponse{
		UserId: int32(id),
	}, nil
}
