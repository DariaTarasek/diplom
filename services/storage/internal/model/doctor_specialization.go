package model

type DoctorSpecialization struct {
	DoctorID         UserID `db:"doctor_id"`
	SpecializationID SpecID `db:"specialization_id"`
}
