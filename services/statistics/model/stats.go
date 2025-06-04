package model

type (
	AllStats struct {
		TotalPatients        int
		TotalVisits          int
		TopServices          []TopService
		DoctorAvgVisit       []DoctorAvgVisit
		DoctorCheckStat      []DoctorCheckStat
		DoctorUniquePatients []DoctorUniquePatients
		AgeGroupStat         []AgeGroupStat
		NewPatientsThisMonth int
		AvgVisitPerPatient   float32
		TotalIncome          float32
		MonthlyIncome        float32
		ClinicAvgCheck       float32
	}
	TopService struct {
		Name  string
		Count int
	}
	DoctorAvgVisit struct {
		Doctor          string
		AvgWeeklyVisits float32
	}
	DoctorCheckStat struct {
		Doctor   string
		AvgCheck float32
	}
	DoctorUniquePatients struct {
		DoctorID       string
		UniquePatients int
	}
	AgeGroupStat struct {
		AgeGroup string
		Percent  float32
	}
)
