package model

type SpecID int

type Specialization struct {
	ID   SpecID `db:"id"`
	Name string `db:"name"`
}
