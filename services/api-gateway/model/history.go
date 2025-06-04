package model

type (
	HistoryVisits struct {
		ID        int    `json:"id"`
		Date      string `json:"date"`
		DoctorID  int    `json:"doctor_id"`
		Doctor    string `json:"doctor"`
		Diagnose  string `json:"diagnose"`
		Treatment string `json:"treatment"`
	}
)
