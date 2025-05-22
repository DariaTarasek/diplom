package model

type VisitID int

type AppointmentMaterial struct {
	ID           int        `db:"id"`
	VisitID      VisitID    `db:"visit_id"`
	MaterialID   MaterialID `db:"material_id"`
	QuantityUsed int        `db:"quantity_used"`
}
