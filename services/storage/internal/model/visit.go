package model

import "time"

type (
	VisitID int
	Visit   struct {
		ID            VisitID       `db:"visit_id"`
		AppointmentID AppointmentID `db:"appointment_id"`
		Complaints    string        `db:"complaints"`
		TreatmentPlan string        `db:"treatment_plan"`
		CreatedAt     time.Time     `db:"created_at"`
		UpdatedAt     time.Time     `db:"updated_at"`
	}
)
