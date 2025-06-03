package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/patient/model"
	authpb "github.com/DariaTarasek/diplom/services/patient/proto/auth"
	storagepb "github.com/DariaTarasek/diplom/services/patient/proto/storage"
	"strings"
)

func (s *PatientService) GetHistoryVisits(ctx context.Context, token string) ([]model.HistoryVisits, error) {
	userID, err := s.AuthClient.Client.GetPatient(ctx, &authpb.GetPatientRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить идентификатор пользователя: %w", err)
	}
	visitsResp, err := s.StorageClient.Client.GetPatientVisits(ctx, &storagepb.GetByIdRequest{Id: userID.Patient.UserId})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить проведенные приемы пользователя: %w", err)
	}
	icdCodesResp, err := s.StorageClient.Client.GetICDCodes(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить МКБ коды: %w", err)
	}
	icdMap := make(map[int]string)
	for _, icd := range icdCodesResp.IcdCode {
		icdMap[int(icd.Id)] = fmt.Sprintf("%s: %s", icd.Code, icd.Name)
	}
	res := make([]model.HistoryVisits, 0, len(visitsResp.Visits))
	for _, visit := range visitsResp.Visits {
		diangosesResp, err := s.StorageClient.Client.GetDiagnoseByVisitID(ctx, &storagepb.GetByIDRequest{Id: visit.Id})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить диагнозы приема: %w", err)
		}
		diagnose := make([]string, 0, len(diangosesResp.Diagnose))
		for _, d := range diangosesResp.Diagnose {
			diagnose = append(diagnose, fmt.Sprintf("%s (%s)", icdMap[int(d.IcdCodeId)], d.Note))
		}
		docResp, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: visit.DoctorId})
		if err != nil {
			return nil, fmt.Errorf("не удалось получить врача, проводившего прием: %w", err)
		}
		historyVisit := model.HistoryVisits{
			ID:        int(visit.Id),
			Date:      visit.CreatedAt.AsTime().Format("02.01.2006"),
			DoctorID:  int(visit.DoctorId),
			Doctor:    fmt.Sprintf("%s %s.%s.", docResp.Doctor.SecondName, getAndCapitalizeFirstLetter(docResp.Doctor.FirstName), getAndCapitalizeFirstLetter(docResp.Doctor.Surname)),
			Diagnose:  strings.Join(diagnose, "; "),
			Treatment: visit.Treatment,
		}
		res = append(res, historyVisit)
	}
	return res, nil
}
