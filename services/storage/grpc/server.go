package grpcserver

import (
	"context"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/DariaTarasek/diplom/services/storage/internal/store"
	pb "github.com/DariaTarasek/diplom/services/storage/proto"
)

type Server struct {
	pb.UnimplementedStorageServiceServer
	Store *store.Store
}

func (s *Server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	id, err := s.Store.AddUser(ctx, model.User{
		Login:    &req.Login,
		Password: &req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddUserResponse{
		UserId: int32(id),
	}, nil
}

func (s *Server) AddDoctor(ctx context.Context, req *pb.AddDoctorRequest) (*pb.AddDoctorResponse, error) {
	exp := int(req.Experience)
	err := s.Store.AddDoctor(ctx, model.Doctor{
		ID:          model.UserID(req.UserId),
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Surname:     &req.Surname,
		PhoneNumber: &req.PhoneNumber,
		Email:       req.Email,
		Education:   &req.Education,
		Experience:  &exp,
		Gender:      req.Gender,
	})
	if err != nil {
		return nil, err
	}

	return &pb.AddDoctorResponse{}, nil
}
