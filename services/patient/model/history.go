package model

type (
	HistoryVisits struct {
		ID        int
		Date      string
		DoctorID  int
		Doctor    string
		Diagnose  string
		Treatment string
	}
)
