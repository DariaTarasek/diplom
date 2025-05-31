package grpcserver

import (
	"context"
	"fmt"
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
func (s *Server) GetDoctorSpecsByDoctorId(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetDoctorSpecsByDoctorIdResponse, error) {
	items, err := s.Store.GetDoctorSpecializations(ctx, model.UserID(req.Id))
	if err != nil {
		return nil, err
	}
	var specs []int32
	for _, item := range items {
		specs = append(specs, int32(item.SpecializationID))
	}
	return &pb.GetDoctorSpecsByDoctorIdResponse{Specs: specs}, nil
}

func (s *Server) GetAdmins(ctx context.Context, req *pb.EmptyRequest) (*pb.GetAdminsResponse, error) {
	items, err := s.Store.GetAdmins(ctx)
	if err != nil {
		return nil, err
	}
	var admins []*pb.Admin
	for _, item := range items {
		admin := &pb.Admin{
			UserId:      int32(item.ID),
			FirstName:   item.FirstName,
			SecondName:  item.SecondName,
			Surname:     deref(item.Surname),
			PhoneNumber: deref(item.PhoneNumber),
			Email:       item.Email,
			Gender:      item.Gender,
		}
		admins = append(admins, admin)
	}
	return &pb.GetAdminsResponse{Admins: admins}, nil
}

func (s *Server) GetPatients(ctx context.Context, req *pb.EmptyRequest) (*pb.GetPatientsResponse, error) {
	items, err := s.Store.GetPatients(ctx)
	if err != nil {
		return nil, err
	}
	var patients []*pb.Patient
	for _, item := range items {
		patient := &pb.Patient{
			UserId:      int32(item.ID),
			FirstName:   item.FirstName,
			Surname:     *item.Surname,
			SecondName:  item.SecondName,
			Email:       *item.Email,
			BirthDate:   timestamppb.New(item.BirthDate),
			PhoneNumber: *item.PhoneNumber,
			Gender:      item.Gender,
		}
		patients = append(patients, patient)
	}
	return &pb.GetPatientsResponse{Patients: patients}, nil
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

func (s *Server) GetDoctorWeeklySchedule(ctx context.Context, req *pb.GetScheduleByDoctorIdRequest) (*pb.GetScheduleByDoctorIdResponse, error) {
	items, err := s.Store.GetScheduleByDoctorID(ctx, model.UserID(req.DoctorId))
	if err != nil {
		return nil, err
	}
	var schedule []*pb.WeeklyDoctorSchedule
	for _, item := range items {
		start := *item.StartTime
		end := *item.EndTime
		day := &pb.WeeklyDoctorSchedule{
			Id:                  int32(item.ID),
			DoctorId:            int32(item.DoctorID),
			Weekday:             int32(item.Weekday),
			StartTime:           timestamppb.New(start),
			EndTime:             timestamppb.New(end),
			SlotDurationMinutes: int32(*item.SlotDurationMinutes),
			IsDayOff:            *item.IsDayOff,
		}
		schedule = append(schedule, day)
	}
	return &pb.GetScheduleByDoctorIdResponse{DoctorSchedule: schedule}, nil
}

func (s *Server) UpdateClinicWeeklySchedule(ctx context.Context, request *pb.UpdateClinicWeeklyScheduleRequest) (*pb.DefaultResponse, error) {
	var schedule []model.ClinicSchedule
	for _, item := range request.ClinicSchedule {
		start := item.StartTime.AsTime()
		end := item.EndTime.AsTime()
		duration := int(item.SlotDurationMinutes)
		day := model.ClinicSchedule{
			ID:                  int(item.Id),
			Weekday:             int(item.Weekday),
			StartTime:           &start,
			EndTime:             &end,
			SlotDurationMinutes: &duration,
			IsDayOff:            &item.IsDayOff,
		}
		schedule = append(schedule, day)
	}
	err := s.Store.UpdateClinicSchedule(ctx, schedule)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateDoctorWeeklySchedule(ctx context.Context, request *pb.UpdateDoctorWeeklyScheduleRequest) (*pb.DefaultResponse, error) {
	var schedule []model.DoctorSchedule
	for _, item := range request.DoctorSchedule {
		start := item.StartTime.AsTime()
		end := item.EndTime.AsTime()
		duration := int(item.SlotDurationMinutes)
		day := model.DoctorSchedule{
			ID:                  int(item.Id),
			DoctorID:            model.UserID(item.DoctorId),
			Weekday:             int(item.Weekday),
			StartTime:           &start,
			EndTime:             &end,
			SlotDurationMinutes: &duration,
			IsDayOff:            &item.IsDayOff,
		}
		schedule = append(schedule, day)
	}
	err := s.Store.UpdateDoctorSchedule(ctx, schedule)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) GetRolePermission(ctx context.Context, req *pb.GetRolePermissionRequest) (*pb.DefaultResponse, error) {
	_, err := s.Store.GetRolePermission(ctx, model.RoleID(req.RoleId), model.PermissionID(req.PermId))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, err
}

func (s *Server) GetDoctorsBySpecID(ctx context.Context, req *pb.GetDoctorBySpecIDRequest) (*pb.GetDoctorsResponse, error) {
	fmt.Println(req.SpecId)
	items, err := s.Store.GetDoctorsBySpecializationID(ctx, model.SpecID(req.SpecId))
	if err != nil {
		return nil, fmt.Errorf("не удалось получить врачей по id специальности: %w", err)
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

func (s *Server) GetMaterials(ctx context.Context, req *pb.EmptyRequest) (*pb.GetMaterialsResponse, error) {
	items, err := s.Store.GetMaterials(ctx)
	if err != nil {
		return nil, err
	}
	var materials []*pb.Material
	for _, item := range items {
		material := &pb.Material{
			Id:    int32(item.ID),
			Name:  item.Name,
			Price: int32(item.Price),
		}
		materials = append(materials, material)
	}
	return &pb.GetMaterialsResponse{Materials: materials}, nil
}

func (s *Server) GetServices(ctx context.Context, req *pb.EmptyRequest) (*pb.GetServicesResponse, error) {
	items, err := s.Store.GetServices(ctx)
	if err != nil {
		return nil, err
	}
	var services []*pb.Service
	for _, item := range items {
		service := &pb.Service{
			Id:    int32(item.ID),
			Name:  item.Name,
			Price: int32(*item.Price),
			Type:  int32(item.Category),
		}
		services = append(services, service)
	}
	return &pb.GetServicesResponse{Services: services}, nil
}

func (s *Server) GetServicesTypes(ctx context.Context, req *pb.EmptyRequest) (*pb.GetServicesTypesResponse, error) {
	items, err := s.Store.GetServiceTypes(ctx)
	if err != nil {
		return nil, err
	}
	var types []*pb.ServiceType
	for _, item := range items {
		serviceType := &pb.ServiceType{
			Id:   int32(item.ID),
			Name: item.Name,
		}
		types = append(types, serviceType)
	}
	return &pb.GetServicesTypesResponse{Types: types}, nil
}

func (s *Server) GetServiceTypeById(ctx context.Context, req *pb.GetServiceTypeByIdRequest) (*pb.GetServiceTypeByIdResponse, error) {
	serviceType, err := s.Store.GetServiceTypeByID(ctx, model.ServiceTypeID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetServiceTypeByIdResponse{
		Id:   int32(serviceType.ID),
		Name: serviceType.Name,
	}, nil
}

func (s *Server) AddMaterial(ctx context.Context, req *pb.AddMaterialRequest) (*pb.DefaultResponse, error) {
	_, err := s.Store.AddMaterial(ctx, model.Material{
		Name:  req.Name,
		Price: int(req.Price),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddService(ctx context.Context, req *pb.AddServiceRequest) (*pb.DefaultResponse, error) {
	price := int(req.Price)
	_, err := s.Store.AddService(ctx, model.Service{
		Name:     req.Name,
		Price:    &price,
		Category: model.ServiceTypeID(req.Type),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateMaterial(ctx context.Context, req *pb.UpdateMaterialRequest) (*pb.DefaultResponse, error) {
	err := s.Store.UpdateMaterial(ctx, model.MaterialID(req.Id), model.Material{
		ID:    model.MaterialID(req.Id),
		Name:  req.Name,
		Price: int(req.Price),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.DefaultResponse, error) {
	price := int(req.Price)
	err := s.Store.UpdateService(ctx, model.ServiceID(req.Id), model.Service{
		ID:       model.ServiceID(req.Id),
		Name:     req.Name,
		Price:    &price,
		Category: model.ServiceTypeID(req.Type),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteMaterial(ctx context.Context, req *pb.DeleteRequest) (*pb.DefaultResponse, error) {
	err := s.Store.DeleteMaterial(ctx, model.MaterialID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteService(ctx context.Context, req *pb.DeleteRequest) (*pb.DefaultResponse, error) {
	err := s.Store.DeleteService(ctx, model.ServiceID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddClinicDailyOverride(ctx context.Context, req *pb.AddClinicDailyOverrideRequest) (*pb.DefaultResponse, error) {
	start := req.StartTime.AsTime()
	end := req.EndTime.AsTime()
	duration := int(req.SlotDurationMinutes)
	isDayOff := req.IsDayOff
	_, err := s.Store.AddClinicOverride(ctx, model.ClinicDailyOverride{
		Date:                req.Date.AsTime(),
		StartTime:           &start,
		EndTime:             &end,
		SlotDurationMinutes: &duration,
		IsDayOff:            &isDayOff,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, err
}

func (s *Server) AddDoctorDailyOverride(ctx context.Context, req *pb.AddDoctorDailyOverrideRequest) (*pb.DefaultResponse, error) {
	start := req.StartTime.AsTime()
	end := req.EndTime.AsTime()
	duration := int(req.SlotDurationMinutes)
	isDayOff := req.IsDayOff
	_, err := s.Store.AddDoctorOverride(ctx, model.DoctorDailyOverride{
		DoctorID:            model.UserID(req.DoctorId),
		Date:                req.Date.AsTime(),
		StartTime:           &start,
		EndTime:             &end,
		SlotDurationMinutes: &duration,
		IsDayOff:            &isDayOff,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, err
}

func (s *Server) GetClinicOverride(ctx context.Context, req *pb.GetClinicOverrideRequest) (*pb.GetClinicOverrideResponse, error) {
	override, err := s.Store.GetClinicOverridesByDate(ctx, req.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetClinicOverrideResponse{
		Date:      timestamppb.New(override.Date),
		StartTime: timestamppb.New(*override.StartTime),
		EndTime:   timestamppb.New(*override.EndTime),
		IsDayOff:  *override.IsDayOff,
	}, nil
}

func (s *Server) GetDoctorOverride(ctx context.Context, req *pb.GetDoctorOverrideRequest) (*pb.GetDoctorOverrideResponse, error) {
	override, err := s.Store.GetOverridesByDoctorAndDate(ctx, model.UserID(req.DoctorId), req.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return &pb.GetDoctorOverrideResponse{
		DoctorId:  int32(override.DoctorID),
		Date:      timestamppb.New(override.Date),
		StartTime: timestamppb.New(*override.StartTime),
		EndTime:   timestamppb.New(*override.EndTime),
		IsDayOff:  *override.IsDayOff,
	}, nil
}

func (s *Server) UpdatePatient(ctx context.Context, req *pb.UpdatePatientRequest) (*pb.DefaultResponse, error) {
	err := s.Store.UpdatePatient(ctx, model.UserID(req.UserId), model.Patient{
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
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateDoctor(ctx context.Context, req *pb.UpdateDoctorRequest) (*pb.DefaultResponse, error) {
	exp := int(req.Experience)
	err := s.Store.UpdateDoctor(ctx, model.UserID(req.UserId), model.Doctor{
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
	return &pb.DefaultResponse{}, nil
}

func (s *Server) AddDoctorSpec(ctx context.Context, req *pb.AddDoctorSpecRequest) (*pb.DefaultResponse, error) {
	err := s.Store.AddDoctorSpecialization(ctx, model.DoctorSpecialization{
		DoctorID:         model.UserID(req.DoctorId),
		SpecializationID: model.SpecID(req.SpecId),
	})
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteDoctorSpec(ctx context.Context, req *pb.DeleteDoctorSpecRequest) (*pb.DefaultResponse, error) {
	err := s.Store.DeleteDoctorSpecialization(ctx, model.UserID(req.DoctorId), model.SpecID(req.SpecId))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateAdmin(ctx context.Context, req *pb.UpdateAdminRequest) (*pb.DefaultResponse, error) {
	err := s.Store.UpdateAdmin(ctx, model.UserID(req.UserId), model.Admin{
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
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateAdminRole(ctx context.Context, req *pb.UpdateAdminRoleRequest) (*pb.DefaultResponse, error) {
	err := s.Store.UpdateUserRole(ctx, model.UserID(req.UserId), model.RoleID(req.RoleId))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*pb.DefaultResponse, error) {
	err := s.Store.DeleteUser(ctx, model.UserID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}

func (s *Server) UpdateUserLogin(ctx context.Context, req *pb.UpdateUserLoginRequest) (*pb.DefaultResponse, error) {
	err := s.Store.UpdateUserLogin(ctx, model.UserID(req.UserId), req.Login)
	if err != nil {
		return nil, err
	}
	return &pb.DefaultResponse{}, nil
}
