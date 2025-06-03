package model

type (
	Top3Services struct {
		Name       string `db:"name"`
		UsageCount int    `db:"usage_count"`
	}
	DoctorAvgVisit struct {
		DoctorID        int     `db:"doctor_id"`
		AvgWeeklyVisits float64 `db:"avg_weekly_visits"`
	}
	DoctorCheckStat struct {
		DoctorID int     `db:"doctor_id"`
		AvgCheck float64 `db:"avg_check"`
	}
	DoctorUniquePatients struct {
		DoctorID       int `db:"doctor_id"`
		UniquePatients int `db:"unique_patients"`
	}
	AgeGroupStat struct {
		AgeGroup int     `db:"age_group"`
		Percent  float64 `db:"percent"`
	}
)
