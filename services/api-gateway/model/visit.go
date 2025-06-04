package model

type Diagnose struct {
	ICDCode string `json:"icd_code"`
	Notes   string `json:"notes"`
}

type Visit struct {
	ID            int        `json:"id"`
	AppointmentID int        `json:"appointment_id"`
	PatientID     int        `json:"patient_id"`
	Doctor        string     `json:"doctor"`
	Complaints    string     `json:"complaints"`
	Treatment     string     `json:"treatment"`
	CreatedAt     string     `json:"created_at"`
	Diagnoses     []Diagnose `json:"diagnoses"`
}

type VisitSaveRequest struct {
	AppointmentID int                  `json:"appointment_id"`
	PatientID     int                  `json:"patient_id"`
	DoctorID      int                  `json:"doctor_id"`
	Complaints    string               `json:"complaints"`
	Treatment     string               `json:"treatment"`
	Services      []VisitServiceInput  `json:"manipulations"`
	Materials     []VisitMaterialInput `json:"materials"`
	ICDCodes      []ICDCodeInput       `json:"icd_codes"`
}

type VisitMaterialInput struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

type VisitServiceInput struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

type ICDCodeInput struct {
	CodeID  int    `json:"code"`
	Comment string `json:"comment"`
}

type MaterialsAndServices struct {
	ID       int
	VisitID  int
	Item     string
	Quantity int
}

type VisitPayment struct {
	VisitID              int                    `json:"visit_id"`
	Doctor               string                 `json:"doctor"`
	Patient              string                 `json:"patient"`
	CreatedAt            string                 `json:"created_at"`
	Price                int                    `json:"price"`
	MaterialsAndServices []MaterialsAndServices `json:"materials_and_services"`
}

type VisitPaymentUpdate struct {
	VisitID int    `json:"visit_id"`
	Price   int    `json:"price"`
	Status  string `json:"status"`
}
