package model

import "time"

type PatientMedicalNote struct {
	ID        int       `db:"id"`
	PatientID UserID    `db:"patient_id"`
	Type      string    `db:"type"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
}
