package model

type (
	AppointmentID    int
	UserID           int
	TodayAppointment struct {
		ID        AppointmentID
		Date      string
		Time      string
		PatientID UserID
		Patient   string
	}
	ScheduleTable struct {
		Dates []string                                        // ["01.06.2025", "02.06.2025", ...]
		Times []string                                        // ["09:00", "09:30", "10:00", ...]
		Table map[string]map[string]*UpcomingAdminAppointment // table[date][time] = запись или nil
	}

	UpcomingAdminAppointment struct {
		ID        AppointmentID
		PatientID UserID
		DoctorID  UserID

		PatientSecondName string
		PatientFirstName  string
		PatientSurname    string

		DoctorSecondName string
		DoctorFirstName  string
		DoctorSurname    string
	}

	ScheduleDay struct {
		Date    string `json:"date"`
		Weekday string `json:"weekday"`
	}

	Person struct {
		ID         UserID `json:"id"`
		SecondName string `json:"second_name"`
		FirstName  string `json:"first_name"`
		Surname    string `json:"surname"`
		BirthDate  string // либо time.Time, но строкой для фронта
		Gender     string
		Phone      string
		Specialty  string // для Doctor
	}

	AdminAppointment struct {
		ID      int
		Doctor  Person `json:"doctor"`
		Patient Person `json:"patient"`
	}

	ScheduleMetadata struct {
		Days      []ScheduleDay `json:"days"`
		TimeSlots []string      `json:"timeSlots"`
	}

	AdminScheduleOverview struct {
		Schedule     ScheduleMetadata                         `json:"schedule"`
		Appointments map[string]map[string][]AdminAppointment `json:"appointments"`
	}
)
