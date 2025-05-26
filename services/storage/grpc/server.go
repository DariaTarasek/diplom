package grpcserver

import (
	"context"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/DariaTarasek/diplom/services/storage/internal/store"
	pb "github.com/DariaTarasek/diplom/services/storage/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedStorageServiceServer
	Store *store.Store
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
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

func (s *Server) AddPatient(ctx context.Context, req *pb.AddPatientRequest) (*pb.AddPatientResponse, error) {
	err := s.Store.AddPatient(ctx, model.Patient{
		ID:          model.UserID(req.UserId),
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Surname:     &req.Surname,
		Email:       &req.Email,
		BirthDate:   req.BirthDate.AsTime(),
		PhoneNumber: &req.PhoneNumber,
		Gender:      req.Gender,
	})

	if err != nil {
		return nil, err
	}
	return &pb.AddPatientResponse{}, nil
}

func (s *Server) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.GetUserByLoginResponse, error) {
	user, err := s.Store.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByLoginResponse{
		Login:    deref(user.Login),
		Password: deref(user.Password),
		Id:       int32(user.ID),
	}, nil
}

func (s *Server) GetUserRole(ctx context.Context, req *pb.GetUserRoleRequest) (*pb.GetUserRoleResponse, error) {
	userRole, err := s.Store.GetRoleByUser(ctx, model.UserID(req.UserId))
	if err != nil {
		return nil, err
	}
	return &pb.GetUserRoleResponse{Role: int32(userRole.RoleID)}, nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*pb.DefaultResponse, error) {
	err := s.Store.UpdateUser(ctx, model.UserID(req.Id), model.User{
		Login:    &req.Login,
		Password: &req.Password,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetDoctors(ctx context.Context, req *pb.EmptyRequest) (*pb.GetDoctorsResponse, error) {
	items, err := s.Store.GetDoctors(ctx)
	if err != nil {
		return nil, err
	}
	var doctors []*pb.Doctor
	for _, item := range items {
		doctor := &pb.Doctor{
			UserId:      int32(item.ID),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     deref(item.Surname),
			PhoneNumber: deref(item.PhoneNumber),
			Email:       item.Email,
			Education:   deref(item.Education),
			Experience:  int32(derefInt(item.Experience)),
			Gender:      item.Gender,
		}
		doctors = append(doctors, doctor)
	}
	return &pb.GetDoctorsResponse{Doctors: doctors}, nil
}

func (s *Server) GetClinicWeeklySchedule(ctx context.Context, req *pb.EmptyRequest) (*pb.GetClinicWeeklyScheduleResponse, error) {
	items, err := s.Store.GetClinicSchedule(ctx)
	if err != nil {
		return nil, err
	}
	var schedule []*pb.WeeklyClinicSchedule
	for _, item := range items {
		start := *item.StartTime
		end := *item.EndTime
		day := &pb.WeeklyClinicSchedule{
			Id:                  int32(item.ID),
			Weekday:             int32(item.Weekday),
			StartTime:           timestamppb.New(start),
			EndTime:             timestamppb.New(end),
			SlotDurationMinutes: int32(*item.SlotDurationMinutes),
			IsDayOff:            *item.IsDayOff,
		}
		schedule = append(schedule, day)
	}
	return &pb.GetClinicWeeklyScheduleResponse{ClinicSchedule: schedule}, nil
}
