package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AdminService) GetUnconfirmedAppointments(ctx context.Context) ([]model.Appointment, error) {
	resp, err := s.StorageClient.Client.GetAppointments(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	var unconfirmedAppts []model.Appointment
	for _, item := range resp.Appointments {
		if item.Status == "unconfirmed" {

			doctorResp, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: item.DoctorId})
			if err != nil {
				return nil, err
			}
			doctorName := fmt.Sprintf("%s %s %s", doctorResp.Doctor.SecondName, doctorResp.Doctor.FirstName, doctorResp.Doctor.Surname)

			appt := model.Appointment{
				ID:                int(item.Id),
				Doctor:            doctorName,
				PatientID:         int(item.PatientId),
				Date:              item.Date.AsTime().Format("02.01.2006"),
				Time:              item.Time.AsTime().Format("15:04"),
				PatientFirstName:  item.FirstName,
				PatientSecondName: item.SecondName,
				PatientSurname:    item.Surname,
				PatientBirthDate:  item.BirthDate.AsTime().Format("02.01.2006"),
				Gender:            item.Gender,
				PhoneNumber:       item.PhoneNumber,
				Status:            item.Status,
				CreatedAt:         item.CreatedAt.AsTime().Format("02.01.2006"),
				UpdatedAt:         item.UpdatedAt.AsTime().Format("02.01.2006"),
			}
			unconfirmedAppts = append(unconfirmedAppts, appt)
		}
	}
	return unconfirmedAppts, nil
}

func (s *AdminService) UpdateAppointment(ctx context.Context, appt model.UpdateAppointment) error {
	apptReq := &storagepb.Appointment{
		Id:        int32(appt.ID),
		Date:      timestamppb.New(appt.Date),
		Time:      timestamppb.New(appt.Time),
		Status:    appt.Status,
		UpdatedAt: timestamppb.New(appt.UpdatedAt),
	}
	_, err := s.StorageClient.Client.UpdateAppointment(ctx, &storagepb.UpdateAppointmentRequest{Appointment: apptReq})
	if err != nil {
		return err
	}
	return nil
}
