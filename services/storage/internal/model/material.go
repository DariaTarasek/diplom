package model

type MaterialID int

type Material struct {
	ID    MaterialID `db:"id"`
	Name  string     `db:"name"`
	Price int        `db:"price"`
}
