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
		ID                 AppointmentID `json:"id"`
		DoctorID           UserID        `json:"doctor_id"`
		PatientID          *UserID       `json:"user_id"`
		Date               string        `json:"date"`
		Time               string        `json:"time"`
		PatientSecondName  string        `json:"secondName"`
		PatientFirstName   string        `json:"firstName"`
		PatientSurname     *string       `json:"surname"`
		PatientBirthDate   string        `json:"birthDate"`
		PatientGender      string        `json:"gender"`
		PatientPhoneNumber string        `json:"phone"`
		Status             string        `json:"status"`
		CreatedAt          string        `json:"createdAt"`
		UpdatedAt          string        `json:"updatedAt"`
	}
	UnconfirmedAppointment struct {
		ID                int    `json:"id"`
		Doctor            string `json:"doctor"`
		PatientID         int    `json:"patient_id"`
		Date              string `json:"date"`
		Time              string `json:"time"`
		PatientFirstName  string `json:"patient_first_name"`
		PatientSecondName string `json:"patient_second_name"`
		PatientSurname    string `json:"patient_surname"`
		PatientBirthDate  string `json:"patient_birth_date"`
		Gender            string `json:"gender"`
		PhoneNumber       string `json:"phone_number"`
		Status            string `json:"status"`
		CreatedAt         string `json:"created_at"`
		UpdatedAt         string `json:"updated_at"`
	}
	UpdateAppointment struct {
		ID        int       `json:"id"`
		Date      string    `json:"date"`
		Time      string    `json:"time"`
		Status    string    `json:"status"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	UpcomingAppointment struct {
		ID        AppointmentID `json:"id"`
		Date      string        `json:"date"`
		Time      string        `json:"time"`
		DoctorID  UserID        `json:"doctorId"`
		Doctor    string        `json:"doctor"`
		Specialty string        `json:"specialty"`
	}
	TodayAppointment struct {
		ID        AppointmentID `json:"id"`
		Date      string        `json:"date"`
		Time      string        `json:"time"`
		PatientID UserID        `json:"patient_id"`
		Patient   string        `json:"patient"`
	}
	ScheduleTable struct {
		Dates []string                                         `json:"dates"` // ["01.06.2025", "02.06.2025", ...]
		Times []string                                         `json:"times"` // ["09:00", "09:30", "10:00", ...]
		Table map[string]map[string]*UpcomingDoctorAppointment `json:"table"` // table[date][time] = запись или nil
	}

	UpcomingDoctorAppointment struct {
		ID        AppointmentID `json:"id"`
		PatientID UserID        `json:"patient_id"`
		Patient   string        `json:"patient"`
		// другие поля при необходимости
	}
)
