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

func (s *Server) GetAllSpecs(ctx context.Context, req *pb.EmptyRequest) (*pb.GetAllSpecsResponse, error) {
	items, err := s.Store.GetSpecializations(ctx)
	if err != nil {
		return nil, err
	}
	var specs []*pb.Specialization
	for _, item := range items {
		spec := &pb.Specialization{
			Id:   int32(item.ID),
			Name: item.Name,
		}
		specs = append(specs, spec)
	}
	return &pb.GetAllSpecsResponse{Specs: specs}, nil
}

func (s *Server) AddUserRole(ctx context.Context, req *pb.AddUserRoleRequest) (*pb.AddUserRoleResponse, error) {
	err := s.Store.AddUserRole(ctx, model.UserID(req.UserId), model.RoleID(req.RoleId))
	if err != nil {
		return nil, err
	}
	return &pb.AddUserRoleResponse{}, nil
}

func (s *Server) AddAdmin(ctx context.Context, req *pb.AddAdminRequest) (*pb.AddAdminResponse, error) {
	err := s.Store.AddAdmin(ctx, model.Admin{
		ID:          model.UserID(req.UserId),
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Surname:     &req.Surname,
		PhoneNumber: &req.PhoneNumber,
		Email:       req.Email,
		Gender:      req.Gender,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddAdminResponse{}, err
}
