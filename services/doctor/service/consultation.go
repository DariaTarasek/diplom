package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/doctor/model"
	authpb "github.com/DariaTarasek/diplom/services/doctor/proto/auth"
	storagepb "github.com/DariaTarasek/diplom/services/doctor/proto/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *DoctorService) GetPatientAllergiesChronics(ctx context.Context, id int) ([]model.AllergiesChronics, error) {
	notesResp, err := s.StorageClient.Client.GetPatientAllergiesChronics(ctx, &storagepb.GetByIdRequest{Id: int32(id)})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список аллергий и хронческих заболеваний: %w", err)
	}
	var notes []model.AllergiesChronics
	for _, item := range notesResp.PatientAllergiesChronics {
		note := model.AllergiesChronics{
			ID:        int(item.Id),
			PatientID: int(item.PatientId),
			Type:      item.Type,
			Title:     item.Title,
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (s *DoctorService) GetAppointmentByID(ctx context.Context, id int) (model.Appointment, error) {
	apptResp, err := s.StorageClient.Client.GetAppointmentByID(ctx, &storagepb.GetByIDRequest{Id: int32(id)})
	if err != nil {
		return model.Appointment{}, fmt.Errorf("не удалось получить запись: %w", err)
	}

	patientId := model.UserID(apptResp.Appointment.PatientId)
	appt := model.Appointment{
		ID:                 model.AppointmentID(apptResp.Appointment.Id),
		DoctorID:           model.UserID(apptResp.Appointment.DoctorId),
		PatientID:          &patientId,
		Date:               apptResp.Appointment.Date.AsTime(),
		Time:               apptResp.Appointment.Time.AsTime(),
		PatientSecondName:  apptResp.Appointment.SecondName,
		PatientFirstName:   apptResp.Appointment.FirstName,
		PatientSurname:     &apptResp.Appointment.Surname,
		PatientBirthDate:   apptResp.Appointment.BirthDate.AsTime(),
		PatientGender:      apptResp.Appointment.Gender,
		PatientPhoneNumber: apptResp.Appointment.PhoneNumber,
		Status:             apptResp.Appointment.Status,
		CreatedAt:          apptResp.Appointment.CreatedAt.AsTime(),
		UpdatedAt:          apptResp.Appointment.UpdatedAt.AsTime(),
	}

	return appt, nil
}

// GetPatientVisits Получение визитов и диагнозов пациента
func (s *DoctorService) GetPatientVisits(ctx context.Context, id int) ([]model.Visit, error) {
	icdResp, err := s.StorageClient.Client.GetICDCodes(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить список МКБ-кодов: %w", err)
	}
	icdCodes := make(map[int]string)
	for _, item := range icdResp.IcdCode {
		icdCodes[int(item.Id)] = fmt.Sprintf("%s: %s", item.Code, item.Name)
	}

	visitsResp, err := s.StorageClient.Client.GetPatientVisits(ctx, &storagepb.GetByIdRequest{Id: int32(id)})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить визиты пациента: %w", err)
	}
	var visits []model.Visit
	for _, item := range visitsResp.Visits {

		/// получение диагнозов, поставленных во время визита
		diagnosesResp, err := s.StorageClient.Client.GetPatientDiagnoses(ctx, &storagepb.GetByIdRequest{Id: item.Id})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить диагнозы по айди визита: %w", err)
		}
		var diagnoses []model.Diagnose
		for _, d := range diagnosesResp.Diagnoses {
			diagnose := model.Diagnose{
				ICDCode: icdCodes[int(d.IcdCodeId)],
				Notes:   d.Note,
			}
			diagnoses = append(diagnoses, diagnose)
		}

		// получение врача, проводившего визит
		doctorResp, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: item.DoctorId})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить врача, проводившего прием: %w", err)
		}
		doctor := fmt.Sprintf("%s %s %s", doctorResp.Doctor.SecondName, doctorResp.Doctor.FirstName, doctorResp.Doctor.Surname)
		visit := model.Visit{
			ID:            int(item.Id),
			AppointmentID: int(item.AppointmentId),
			PatientID:     int(item.PatientId),
			Doctor:        doctor,
			Complaints:    item.Complaints,
			Treatment:     item.Treatment,
			CreatedAt:     item.CreatedAt.AsTime().Format("02.01.2006"),
			Diagnoses:     diagnoses,
		}
		visits = append(visits, visit)
	}
	return visits, nil
}
func (s *DoctorService) AddConsultation(ctx context.Context, materials []model.VisitMaterial,
	services []model.VisitService, diagnoses []model.VisitDiagnose, visit model.AddVisit, token string) error {
	apptID := visit.AppointmentID
	appointment, err := s.StorageClient.Client.GetAppointmentByID(ctx, &storagepb.GetByIDRequest{Id: int32(apptID)})
	appointmentReq := appointment.Appointment
	appointmentReq.Status = "completed"
	fmt.Println(appointmentReq)
	appointmentReq.UpdatedAt = timestamppb.New(time.Now())
	_, err = s.StorageClient.Client.UpdateAppointment(ctx, &storagepb.UpdateAppointmentRequest{Appointment: appointmentReq})

	doctorID, err := s.AuthClient.Client.GetUserID(ctx, &authpb.GetUserIDRequest{Token: token})
	if err != nil {
		return fmt.Errorf("не удалось получить доктора, проводившего прием: %w", err)
	}

	resp, err := s.StorageClient.Client.AddPatientVisit(ctx, &storagepb.AddPatientVisitRequest{
		AppointmentId: int32(visit.AppointmentID),
		PatientId:     int32(visit.PatientID),
		DoctorId:      doctorID.UserId,
		Complaints:    visit.Complaints,
		Treatment:     visit.Treatment,
	})
	if err != nil {
		return fmt.Errorf("не удалось добавить визит пациента: %w", err)
	}
	visitID := resp.Id

	var diagnosesReq []*storagepb.Diagnose
	for _, item := range diagnoses {
		req := &storagepb.Diagnose{
			VisitId:   visitID,
			IcdCodeId: int32(item.ICDCodeID),
			Note:      item.Note,
		}
		diagnosesReq = append(diagnosesReq, req)
	}
	_, err = s.StorageClient.Client.AddPatientDiagnoses(ctx, &storagepb.AddPatientDiagnosesRequest{Diagnoses: diagnosesReq})
	if err != nil {
		return fmt.Errorf("не удалось добавить диагнозы: %w", err)
	}

	var servicesReq []*storagepb.AddVisitServices
	for _, item := range services {
		req := &storagepb.AddVisitServices{
			VisitId:   visitID,
			ServiceId: int32(item.ServiceID),
			Amount:    int32(item.Amount),
		}
		servicesReq = append(servicesReq, req)
	}
	_, err = s.StorageClient.Client.AddVisitServices(ctx, &storagepb.AddVisitServicesRequest{Services: servicesReq})
	if err != nil {
		return fmt.Errorf("не удалось добавить проведенные на приеме услуги: %w", err)
	}

	var materialsReq []*storagepb.AddVisitMaterials
	for _, item := range materials {
		req := &storagepb.AddVisitMaterials{
			VisitId:    visitID,
			MaterialId: int32(item.MaterialID),
			Amount:     int32(item.Amount),
		}
		materialsReq = append(materialsReq, req)
	}
	_, err = s.StorageClient.Client.AddVisitMaterials(ctx, &storagepb.AddVisitMaterialsRequest{Materials: materialsReq})
	if err != nil {
		return fmt.Errorf("не удалось добавить затраченные на приеме материалы: %w", err)
	}

	totalResp, err := s.StorageClient.Client.CalculateVisitTotal(ctx, &storagepb.CalculateVisitTotalRequest{
		VisitId: visitID,
	})
	if err != nil {
		return fmt.Errorf("не удалось рассчитать итоговую сумму приёма: %w", err)
	}

	_, err = s.StorageClient.Client.AddOrUpdateVisitPayment(ctx, &storagepb.AddOrUpdateVisitPaymentRequest{
		Payment: &storagepb.VisitPayment{
			VisitId: visitID,
			Price:   totalResp.Total,
			Status:  "unconfirmed",
		},
	})
	if err != nil {
		return fmt.Errorf("не удалось сохранить платёж по приёму: %w", err)
	}
	return nil
}

func (s *DoctorService) AddPatientAllergiesChronics(ctx context.Context, notes []model.AllergiesChronics) error {
	var notesReq []*storagepb.PatientAllergiesChronics
	for _, item := range notes {
		noteReq := &storagepb.PatientAllergiesChronics{
			PatientId: int32(item.PatientID),
			Type:      item.Type,
			Title:     item.Title,
		}
		notesReq = append(notesReq, noteReq)
	}
	_, err := s.StorageClient.Client.AddPatientAllergiesChronics(ctx, &storagepb.AddPatientAllergiesChronicsRequest{Notes: notesReq})
	if err != nil {
		return fmt.Errorf("не удалось добавить мед. запись об аллергиях и хр. заболеваниях: %w", err)
	}
	return nil
}

func (s *DoctorService) AddVisitPayment(ctx context.Context, payment model.VisitPayment) error {
	_, err := s.StorageClient.Client.AddVisitPayment(ctx, &storagepb.VisitPaymentRequest{
		VisitId: int32(payment.VisitID),
		Price:   int32(payment.Price),
		Status:  payment.Status,
	})
	if err != nil {
		return fmt.Errorf("не удалось добавить платеж визита: %w", err)
	}
	return nil
}

func (s *DoctorService) UpdateVisitPayment(ctx context.Context, payment model.VisitPayment) error {
	_, err := s.StorageClient.Client.UpdateVisitPayment(ctx, &storagepb.VisitPaymentRequest{
		VisitId: int32(payment.VisitID),
		Price:   int32(payment.Price),
		Status:  payment.Status,
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить платеж визита: %w", err)
	}
	return nil
}
