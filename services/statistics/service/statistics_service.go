package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/statistics/clients"
	"github.com/DariaTarasek/diplom/services/statistics/model"
	storagepb "github.com/DariaTarasek/diplom/services/statistics/proto/storage"
)

type StatisticsService struct {
	StorageClient *clients.StorageClient
}

func NewStatisticsService(client *clients.StorageClient) *StatisticsService {
	return &StatisticsService{
		StorageClient: client,
	}
}

func (s *StatisticsService) GetAllStats(ctx context.Context) (model.AllStats, error) {
	var AllStats model.AllStats // вся статистика

	// получение числа всех пациентов
	totalPatientsResp, err := s.StorageClient.Client.GetTotalPatients(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить количество всех пациентов: %w", err)
	}
	totalPatients := totalPatientsResp.Int

	// получение числа всех проведенных приемов
	totalVisitsResp, err := s.StorageClient.Client.GetTotalVisits(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить количество всех проведенных приемов: %w", err)
	}
	totalVisits := totalVisitsResp.Int

	// получение топ-3 услуг
	topServicesResp, err := s.StorageClient.Client.GetTopServices(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить топ услуг: %w", err)
	}
	var topServices []model.TopService
	for _, item := range topServicesResp.ServiceStats {
		service := model.TopService{
			Name:  item.Name,
			Count: int(item.UsageCount),
		}
		topServices = append(topServices, service)
	}

	// получение ср. количества визитов к врачу в неделю
	avgVisitsResp, err := s.StorageClient.Client.GetDoctorAvgVisit(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить  ср. количество визитов к врачу в неделю: %w", err)
	}
	var avgVisits []model.DoctorAvgVisit
	for _, item := range avgVisitsResp.Visits {
		doctor, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: item.DoctorId})
		if err != nil {
			return model.AllStats{}, err
		}
		visit := model.DoctorAvgVisit{
			Doctor:          StringDoctorName(doctor),
			AvgWeeklyVisits: item.AvgWeeklyVisits,
		}
		avgVisits = append(avgVisits, visit)
	}

	// получение ср. чека врача
	avgDocCheckResp, err := s.StorageClient.Client.GetDoctorAvgCheck(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить  ср. чек врача: %w", err)
	}
	var avgDocCheck []model.DoctorCheckStat
	for _, item := range avgDocCheckResp.Check {
		doctor, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: item.DoctorId})
		if err != nil {
			return model.AllStats{}, err
		}
		check := model.DoctorCheckStat{
			Doctor:   StringDoctorName(doctor),
			AvgCheck: item.AvgCheck,
		}
		avgDocCheck = append(avgDocCheck, check)
	}

	// получение кол-ва пациентов врача
	docPatientsResp, err := s.StorageClient.Client.GetDoctorUniquePatient(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить кол-во пациентов врача: %w", err)
	}
	var docPatients []model.DoctorUniquePatients
	for _, item := range docPatientsResp.Patients {
		doctor, err := s.StorageClient.Client.GetDoctorByID(ctx, &storagepb.GetByIDRequest{Id: item.DoctorId})
		if err != nil {
			return model.AllStats{}, err
		}
		pat := model.DoctorUniquePatients{
			DoctorID:       StringDoctorName(doctor),
			UniquePatients: int(item.UniquePatients),
		}
		docPatients = append(docPatients, pat)
	}

	// получение возрастных категорий пациентов
	ageGroupsResp, err := s.StorageClient.Client.GetAgeGroupStat(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить возрастные группы пациентов: %w", err)
	}
	var ageGroups []model.AgeGroupStat
	for _, item := range ageGroupsResp.AgeGroups {
		group := model.AgeGroupStat{
			AgeGroup: fmt.Sprintf("%d-%d", item.AgeGroup, item.AgeGroup+9),
			Percent:  item.Percent,
		}
		ageGroups = append(ageGroups, group)
	}

	// получение кол-ва новых пациентов в текущем месяце
	newPatientsThisMonthResp, err := s.StorageClient.Client.GetNewPatientsThisMonth(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить количество новых пациентов: %w", err)
	}
	newPatientsThisMonth := newPatientsThisMonthResp.Int

	// получение ср. кол-ва визитов на одного пациента
	avgVisitsPerPatientResp, err := s.StorageClient.Client.GetAvgVisitsPerPatient(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить количество ср. кол-во визитов на одного пациента: %w", err)
	}
	avgVisitsPerPatient := avgVisitsPerPatientResp.Float

	// получение общей выручки
	totalIncomeResp, err := s.StorageClient.Client.GetTotalIncome(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить общую выручку клиники: %w", err)
	}
	totalIncome := totalIncomeResp.Float

	// получение выручки за текущий месяц
	incomeThisMonthResp, err := s.StorageClient.Client.GetMonthlyIncome(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить выручку за текущий месяц: %w", err)
	}
	incomeThisMonth := incomeThisMonthResp.Float

	// получение среднего чека клиники
	avgCheckResp, err := s.StorageClient.Client.GetClinicAverageCheck(ctx, &storagepb.EmptyRequest{})
	if err != nil {
		return model.AllStats{}, fmt.Errorf("не удалось получить ср. чек клиники: %w", err)
	}
	avgCheck := avgCheckResp.Float

	// заполнение всей статистики
	AllStats.TotalPatients = int(totalPatients)
	AllStats.TotalVisits = int(totalVisits)
	AllStats.TopServices = topServices
	AllStats.DoctorCheckStat = avgDocCheck
	AllStats.DoctorAvgVisit = avgVisits
	AllStats.DoctorUniquePatients = docPatients
	AllStats.AgeGroupStat = ageGroups
	AllStats.NewPatientsThisMonth = int(newPatientsThisMonth)
	AllStats.AvgVisitPerPatient = avgVisitsPerPatient
	AllStats.TotalIncome = totalIncome
	AllStats.MonthlyIncome = incomeThisMonth
	AllStats.ClinicAvgCheck = avgCheck

	return AllStats, nil
}

func StringDoctorName(resp *storagepb.GetDoctorResponse) string {
	return fmt.Sprintf("%s %s %s", resp.Doctor.SecondName, resp.Doctor.FirstName, resp.Doctor.Surname)
}
