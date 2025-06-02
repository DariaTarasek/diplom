package model

import "time"

type (
	VisitID int
	Visit   struct {
		ID            VisitID       `db:"id"`
		AppointmentID AppointmentID `db:"appointment_id"`
		PatientID     UserID        `db:"patient_id"`
		DoctorID      UserID        `db:"doctor_id"`
		Complaints    string        `db:"complaints"`
		TreatmentPlan string        `db:"treatment_plan"`
		CreatedAt     time.Time     `db:"created_at"`
		UpdatedAt     time.Time     `db:"updated_at"`
	}
	VisitPayment struct {
		VisitID VisitID `db:"visit_id"`
		Price   int32   `db:"price"`
		Status  string  `db:"status"`
	}
)
