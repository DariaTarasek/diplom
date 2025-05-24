package model

type AppointmentService struct {
	ID        int       `db:"id"`
	VisitID   VisitID   `db:"visit_id"`
	ServiceID ServiceID `db:"service_id"`
	Quantity  int       `db:"quantity"`
}
