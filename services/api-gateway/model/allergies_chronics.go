package model

type AllergiesChronics struct {
	ID        int    `json:"id"`
	PatientID int    `json:"patient_id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
}
