package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"github.com/DariaTarasek/diplom/services/admin/sharederrors"
)

func (s *AdminService) GetUnconfirmedVisitsPayments(ctx context.Context) ([]model.UnconfirmedVisitPayment, error) {

	// получение счетов на визиты со статусом unconfirmed
	resp, err := s.StorageClient.Client.GetVisitsPayments(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, err
	}
	var visitPayments []model.VisitPayment
	for _, item := range resp.VisitPayment {
		vp := model.VisitPayment{
			VisitID: int(item.VisitId),
			Price:   int(item.Price),
			Status:  item.Status,
		}
		visitPayments = append(visitPayments, vp)
	}

	// получение визитов
	var unconfirmedVisitPayments []model.UnconfirmedVisitPayment
	for _, item := range visitPayments {
		visitsResp, err := s.StorageClient.Client.GetVisitByID(ctx, &storagepb.GetByIdRequest{Id: int32(item.VisitID)})
		if err != nil {
			return nil, err
		}

		// получение ФИО врача по айди
		doctorResp, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: visitsResp.Visit.DoctorId})
		if err != nil {
			return nil, err
		}
		doctor := doctorResp.Doctor
		doctorName := fmt.Sprintf("%s %s %s", doctor.SecondName, doctor.FirstName, doctor.Surname)

		// получение ФИО пациента по айди
		patientResp, err := s.StorageClient.Client.GetPatientByID(ctx, &storagepb.GetByIDRequest{Id: visitsResp.Visit.PatientId})
		if err != nil {
			return nil, err
		}
		patient := patientResp.Patient
		patientName := fmt.Sprintf("%s %s %s", patient.SecondName, patient.FirstName, patient.Surname)

		// заполнение структуры неподтвержденной оплаты
		unconfVP := model.UnconfirmedVisitPayment{
			VisitID:   item.VisitID,
			Doctor:    doctorName,
			Patient:   patientName,
			CreatedAt: visitsResp.Visit.CreatedAt.AsTime().Format("02.01.2006 15:04"),
			Price:     item.Price,
		}
		unconfirmedVisitPayments = append(unconfirmedVisitPayments, unconfVP)
	}
	return unconfirmedVisitPayments, nil
}

func (s *AdminService) UpdateVisitPayment(ctx context.Context, payment model.VisitPayment) error {
	price := payment.Price
	if price < 0 || price > 1000000 {
		return fmt.Errorf("некорректное значение стоимости: %w", sharederrors.ErrInvalidValue)
	}
	paymentReq := &storagepb.VisitPayment{
		VisitId: int32(payment.VisitID),
		Price:   int32(payment.Price),
		Status:  payment.Status,
	}
	_, err := s.StorageClient.Client.AddOrUpdateVisitPayment(ctx, &storagepb.AddOrUpdateVisitPaymentRequest{Payment: paymentReq})
	if err != nil {
		return fmt.Errorf("не удалось подтвердить стоимость приема: %w", err)
	}
	return nil
}
