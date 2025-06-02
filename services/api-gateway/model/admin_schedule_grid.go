package model

type (
	ScheduleDay struct {
		Date    string `json:"date"`
		Weekday string `json:"weekday"`
	}

	Person struct {
		ID         UserID `json:"id"`
		SecondName string `json:"second_name"`
		FirstName  string `json:"first_name"`
		Surname    string `json:"surname"`
		BirthDate  string `json:"birthDate"` // либо time.Time, но строкой для фронта
		Gender     string `json:"gender"`
		Phone      string `json:"phone"`
		Specialty  string `json:"specialty"` // для Doctor
	}

	AdminAppointment struct {
		ID      int    `json:"id"`
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
