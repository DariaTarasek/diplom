package model

import "time"

type (
	UserID        int
	AppointmentID int
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
	TodayAppointment struct {
		ID        AppointmentID
		Date      string
		Time      string
		PatientID UserID
		Patient   string
	}
	ScheduleTable struct {
		Dates []string                                   // ["01.06.2025", "02.06.2025", ...]
		Times []string                                   // ["09:00", "09:30", "10:00", ...]
		Table map[string]map[string]*UpcomingAppointment // table[date][time] = запись или nil
	}

	UpcomingAppointment struct {
		ID        AppointmentID
		PatientID UserID
		Patient   string
		// другие поля при необходимости
	}
)
