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

func (s *Server) AddMaterial(ctx context.Context, req *pb.AddMaterialRequest) (*pb.DefaultResponse, error) {
	err := s.Service.AddMaterial(ctx, model.Material{
		Name:  req.Name,
		Price: int(req.Price),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось добавить материал: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddService(ctx context.Context, req *pb.AddServiceRequest) (*pb.DefaultResponse, error) {
	err := s.Service.AddService(ctx, model.Service{
		Name:     req.Name,
		Price:    int(req.Price),
		Category: int(req.Type),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось добавить услугу: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateMaterial(ctx context.Context, req *pb.UpdateMaterialRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdateMaterial(ctx, model.Material{
		ID:    int(req.Id),
		Name:  req.Name,
		Price: int(req.Price),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить материал: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdateService(ctx, model.Service{
		ID:       int(req.Id),
		Name:     req.Name,
		Price:    int(req.Price),
		Category: int(req.Type),
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить услугу: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteMaterial(ctx context.Context, req *pb.DeleteRequest) (*pb.DefaultResponse, error) {
	err := s.Service.DeleteMaterial(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("не удалось удалить материал: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteService(ctx context.Context, req *pb.DeleteRequest) (*pb.DefaultResponse, error) {
	err := s.Service.DeleteService(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("не удалось удалить услугу: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetAdmins(ctx context.Context, req *pb.EmptyRequest) (*pb.GetAdminsResponse, error) {
	admins, err := s.Service.GetAdmins(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список администраторов: %w", err)
	}
	var gRPCAdmins []*pb.Admin
	for _, item := range admins {
		admin := pb.Admin{
			UserId:      int32(item.ID),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Gender:      item.Gender,
			Role:        item.Role,
		}
		gRPCAdmins = append(gRPCAdmins, &admin)
	}
	return &pb.GetAdminsResponse{Admins: gRPCAdmins}, nil
}
func (s *Server) GetSpecs(ctx context.Context, req *pb.EmptyRequest) (*pb.GetSpecsResponse, error) {
	specs, err := s.Service.GetSpecs(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список специализаций врачей: %w", err)
	}
	var gRPCSpecs []*pb.Spec
	for _, item := range specs {
		spec := pb.Spec{
			Id:   int32(item.ID),
			Name: item.Name,
		}
		gRPCSpecs = append(gRPCSpecs, &spec)
	}
	return &pb.GetSpecsResponse{
		Specs: gRPCSpecs,
	}, nil
}

func (s *Server) GetDoctors(ctx context.Context, req *pb.EmptyRequest) (*pb.GetDoctorsResponse, error) {
	doctors, err := s.Service.GetDoctors(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список врачей: %w", err)
	}
	var gRPCDoctors []*pb.DoctorWithSpecs
	for _, item := range doctors {
		var specs []int32
		for _, specItem := range item.Specs {
			intSpecItem := int32(specItem)
			specs = append(specs, intSpecItem)
		}
		doctor := pb.DoctorWithSpecs{
			UserId:      int32(item.ID),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     item.Surname,
			PhoneNumber: item.PhoneNumber,
			Email:       item.Email,
			Education:   item.Education,
			Experience:  int32(item.Experience),
			Gender:      item.Gender,
			Specs:       specs,
		}
		gRPCDoctors = append(gRPCDoctors, &doctor)
	}
	return &pb.GetDoctorsResponse{Doctors: gRPCDoctors}, nil
}

func (s *Server) GetPatients(ctx context.Context, req *pb.EmptyRequest) (*pb.GetPatientsResponse, error) {
	specs, err := s.Service.GetPatients(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список пациентов: %w", err)
	}
	var gRPCSPatients []*pb.Patient
	for _, item := range specs {
		patient := pb.Patient{
			UserId:      int32(item.ID),
			FirstName:   item.FirstName,
			Surname:     item.Surname,
			SecondName:  item.SecondName,
			Email:       item.Email,
			BirthDate:   item.BirthDate,
			PhoneNumber: item.PhoneNumber,
			Gender:      item.Gender,
		}
		gRPCSPatients = append(gRPCSPatients, &patient)
	}
	return &pb.GetPatientsResponse{Patients: gRPCSPatients}, nil
}

func (s *Server) UpdatePatient(ctx context.Context, req *pb.UpdatePatientRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdatePatient(ctx, model.Patient{
		ID:          int(req.UserId),
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Gender:      req.Gender,
		BirthDate:   req.BirthDate,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить данные пациента: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateDoctor(ctx context.Context, req *pb.UpdateDoctorRequest) (*pb.DefaultResponse, error) {
	var specsInt []int
	for _, item := range req.Specs {
		specsInt = append(specsInt, int(item))
	}
	err := s.Service.UpdateDoctor(ctx, model.Doctor{
		ID:          int(req.UserId),
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Education:   req.Education,
		Experience:  int(req.Experience),
		Gender:      req.Gender,
		Specs:       specsInt,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить данные врача: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateAdmin(ctx context.Context, req *pb.UpdateAdminRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdateAdmin(ctx, model.Admin{
		ID:          int(req.UserId),
		FirstName:   req.FirstName,
		SecondName:  req.SecondName,
		Surname:     req.Surname,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Gender:      req.Gender,
		Role:        req.Role,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить данные администратора: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*pb.DefaultResponse, error) {
	err := s.Service.DeleteUser(ctx, int(req.Id))
	if err != nil {
		return nil, fmt.Errorf("не удалось удалить пользователя: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateEmployeeLogin(ctx context.Context, req *pb.UpdateUserLoginRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdateEmployeeLogin(ctx, int(req.UserId), req.Login)
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить логин сотрудника: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdatePatientLogin(ctx context.Context, req *pb.UpdateUserLoginRequest) (*pb.DefaultResponse, error) {
	err := s.Service.UpdatePatientLogin(ctx, int(req.UserId), req.Login)
	if err != nil {
		return nil, fmt.Errorf("не удалось обновить логин пациента: %w", err)
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetUnconfirmedVisitPayments(ctx context.Context, req *pb.EmptyRequest) (*pb.UnconfirmedVisitPaymentsResponse, error) {
	visits, err := s.Service.GetUnconfirmedVisitsPayments(ctx)
	if err != nil {
		return &pb.UnconfirmedVisitPaymentsResponse{}, fmt.Errorf("не удалось получить визиты с неподтвержденной суммой к оплате: %w", err)
	}
	var visitsResp []*pb.UnconfirmedVisitPayment
	for _, item := range visits {
		visit := &pb.UnconfirmedVisitPayment{
			VisitId:   int32(item.VisitID),
			Doctor:    item.Doctor,
			Patient:   item.Patient,
			CreatedAt: item.CreatedAt,
			Price:     int32(item.Price),
		}
		visitsResp = append(visitsResp, visit)
	}
	return &pb.UnconfirmedVisitPaymentsResponse{VisitPayments: visitsResp}, nil
}
