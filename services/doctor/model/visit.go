package model

type Diagnose struct {
	ICDCode string
	Notes   string
}

type Visit struct {
	ID            int
	AppointmentID int
	PatientID     int
	Doctor        string
	Complaints    string
	Treatment     string
	UpdatedAt     string
	CreatedAt     string
	Diagnoses     []Diagnose
}

type AddVisit struct {
	ID            int
	AppointmentID int
	PatientID     int
	DoctorID      int
	Complaints    string
	Treatment     string
}

type VisitService struct {
	ID        int
	VisitID   int
	ServiceID int
	Amount    int
}

type VisitMaterial struct {
	ID         int
	VisitID    int
	MaterialID int
	Amount     int
}

type VisitDiagnose struct {
	ID        int
	VisitID   int
	ICDCodeID int
	Note      string
}

type VisitPayment struct {
	VisitID int
	Price   int
	Status  string
}
