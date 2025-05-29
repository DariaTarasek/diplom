package model

import "time"

type (
	ScheduleEntry struct {
		Label string   `json:"label"`
		Slots []string `json:"slots"`
	}

	AppointmentID int
	UserID        int
	Appointment   struct {
		ID                 AppointmentID
		DoctorID           UserID
		PatientID          *UserID
		Date               time.Time
		Time               time.Time
		PatientSecondName  string
		PatientFirstName   string
		PatientSurname     *string
		PatientBirthDate   time.Time
		PatientGender      string
		PatientPhoneNumber string
		Status             string
		CreatedAt          time.Time
		UpdatedAt          time.Time
	}
)
