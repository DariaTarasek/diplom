package model

type (
	TopService struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	DoctorAvgVisit struct {
		Doctor          string  `json:"doctor"`
		AvgWeeklyVisits float32 `json:"avgWeeklyVisits"`
	}

	DoctorCheckStat struct {
		Doctor   string  `json:"doctor"`
		AvgCheck float32 `json:"avgCheck"`
	}

	DoctorUniquePatients struct {
		DoctorID       string `json:"doctorId"`
		UniquePatients int    `json:"uniquePatients"`
	}

	AgeGroupStat struct {
		AgeGroup string  `json:"ageGroup"`
		Percent  float32 `json:"percent"`
	}

	AllStats struct {
		TotalPatients        int                    `json:"totalPatients"`
		TotalVisits          int                    `json:"totalVisits"`
		TopServices          []TopService           `json:"topServices"`
		DoctorAvgVisit       []DoctorAvgVisit       `json:"doctorAvgVisit"`
		DoctorCheckStat      []DoctorCheckStat      `json:"doctorCheckStat"`
		DoctorUniquePatients []DoctorUniquePatients `json:"doctorUniquePatients"`
		AgeGroupStat         []AgeGroupStat         `json:"ageGroupStat"`
		NewPatientsThisMonth int                    `json:"newPatientsThisMonth"`
		AvgVisitPerPatient   float32                `json:"avgVisitPerPatient"`
		TotalIncome          float32                `json:"totalIncome"`
		MonthlyIncome        float32                `json:"monthlyIncome"`
		ClinicAvgCheck       float32                `json:"clinicAvgCheck"`
	}
)
