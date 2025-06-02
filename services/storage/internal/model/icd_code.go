package model

type ICDCode struct {
	ID   ICDCodeID `db:"id"`
	Code string    `db:"code"`
	Name string    `db:"name"`
}
