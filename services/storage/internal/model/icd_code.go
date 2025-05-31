package model

type ICDCodeID int

type ICDCode struct {
	ID          ICDCodeID `db:"id"`
	Code        string    `db:"code"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}
