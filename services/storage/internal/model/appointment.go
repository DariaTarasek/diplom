package model

import "time"

type AppointmentID int

type Appointment struct {
	ID                 AppointmentID `db:"id"`
	DoctorID           UserID        `db:"doctor_id"`
	PatientID          *UserID       `db:"patient_id"`
	Date               time.Time     `db:"date"`
	Time               time.Time     `db:"time"`
	PatientSecondName  string        `db:"second_name"`
	PatientFirstName   string        `db:"first_name"`
	PatientSurname     *string       `db:"surname"`
	PatientBirthDate   time.Time     `db:"birth_date"`
	PatientGender      string        `db:"gender"`
	PatientPhoneNumber string        `db:"phone_number"`
	Status             string        `db:"status"`
	CreatedAt          time.Time     `db:"created_at"`
	UpdatedAt          time.Time     `db:"updated_at"`
}
