package model

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
)
