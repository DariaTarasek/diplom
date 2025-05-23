package model

import "time"

type Patient struct {
	ID          UserID
	FirstName   string
	SecondName  string
	Surname     *string
	Email       *string
	BirthDate   time.Time
	PhoneNumber *string
	Gender      string
}
