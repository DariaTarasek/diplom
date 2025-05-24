package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/auth/model"
	"github.com/DariaTarasek/diplom/services/auth/proto/storage"
	"github.com/DariaTarasek/diplom/services/auth/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func (s *AuthService) DoctorRegister(ctx context.Context, user model.User, doctor model.Doctor) (int, error) {
	var hashedPassword string
	var err error
	password := utils.GeneratePassword(12)
	hashedPassword, err = utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("не удалось захешировать пароль: %w", err)
	}

	reqUser := &storagepb.AddUserRequest{
		Login:    *user.Login,
		Password: hashedPassword,
	}

	respUser, err := s.StorageClient.Client.AddUser(ctx, reqUser)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить пользователя через gRPC: %w", err)
	}

	reqDoctor := &storagepb.AddDoctorRequest{
		UserId:      respUser.UserId,
		FirstName:   doctor.FirstName,
		SecondName:  doctor.SecondName,
		Surname:     *doctor.Surname,
		PhoneNumber: *doctor.PhoneNumber,
		Email:       doctor.Email,
		Education:   *doctor.Education,
		Experience:  int32(*doctor.Experience),
		Gender:      doctor.Gender,
	}

	_, err = s.StorageClient.Client.AddDoctor(ctx, reqDoctor)
	if err != nil {
		log.Println(err.Error())
		return 0, fmt.Errorf("не удалось добавить врача через gRPC: %w", err)
	}

	reqUserRole := &storagepb.AddUserRoleRequest{
		UserId: respUser.UserId,
		RoleId: model.DoctorRole,
	}
	_, err = s.StorageClient.Client.AddUserRole(ctx, reqUserRole)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить роль врачу через gRPC: %w", err)
	}

	err = utils.SendPassword(doctor.Email, password)
	if err != nil {
		return 0, fmt.Errorf("не удалось отправить пароль на email: %w", err)
	}

	return int(respUser.UserId), nil
}

func (s *AuthService) AdminRegister(ctx context.Context, user model.User, admin model.Admin, isSuperAdmin bool) (int, error) {
	var hashedPassword string
	var err error
	password := utils.GeneratePassword(12)
	hashedPassword, err = utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("не удалось захешировать пароль: %w", err)
	}

	reqUser := &storagepb.AddUserRequest{
		Login:    *user.Login,
		Password: hashedPassword,
	}

	respUser, err := s.StorageClient.Client.AddUser(ctx, reqUser)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить пользователя через gRPC: %w", err)
	}

	reqAdmin := &storagepb.AddAdminRequest{
		UserId:      respUser.UserId,
		FirstName:   admin.FirstName,
		SecondName:  admin.SecondName,
		Surname:     *admin.Surname,
		PhoneNumber: *admin.PhoneNumber,
		Email:       admin.Email,
		Gender:      admin.Gender,
	}

	_, err = s.StorageClient.Client.AddAdmin(ctx, reqAdmin)
	if err != nil {
		log.Println(err.Error())
		return 0, fmt.Errorf("не удалось добавить администратора через gRPC: %w", err)
	}
	var role model.RoleID
	if isSuperAdmin {
		role = model.SuperAdminRole
	} else {
		role = model.AdminRole
	}
	reqUserRole := &storagepb.AddUserRoleRequest{
		UserId: respUser.UserId,
		RoleId: int32(role),
	}
	_, err = s.StorageClient.Client.AddUserRole(ctx, reqUserRole)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить роль админу через gRPC: %w", err)
	}

	err = utils.SendPassword(admin.Email, password)
	if err != nil {
		return 0, fmt.Errorf("не удалось отправить пароль на email: %w", err)
	}
	return int(respUser.UserId), nil
}

func (s *AuthService) PatientRegister(ctx context.Context, user model.User, patient model.Patient) (int, error) {
	var hashedPassword string
	var err error
	if user.Password != nil {
		hashedPassword, err = utils.HashPassword(*user.Password)
		if err != nil {
			return 0, fmt.Errorf("не удалось захешировать пароль: %w", err)
		}
	}

	reqUser := &storagepb.AddUserRequest{
		Login:    *user.Login,
		Password: hashedPassword,
	}

	respUser, err := s.StorageClient.Client.AddUser(ctx, reqUser)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить пользователя через gRPC: %w", err)
	}

	reqPatient := &storagepb.AddPatientRequest{
		UserId:      respUser.UserId,
		FirstName:   patient.FirstName,
		SecondName:  patient.SecondName,
		Surname:     *patient.Surname,
		BirthDate:   timestamppb.New(patient.BirthDate),
		PhoneNumber: *patient.PhoneNumber,
		Email:       *patient.Email,
		Gender:      patient.Gender,
	}

	_, err = s.StorageClient.Client.AddPatient(ctx, reqPatient)
	if err != nil {
		log.Println(err.Error())
		return 0, fmt.Errorf("не удалось добавить пациента через gRPC: %w", err)
	}

	reqUserRole := &storagepb.AddUserRoleRequest{
		UserId: respUser.UserId,
		RoleId: model.PatientRole,
	}
	_, err = s.StorageClient.Client.AddUserRole(ctx, reqUserRole)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить роль врачу через gRPC: %w", err)
	}

	return int(respUser.UserId), nil
}
