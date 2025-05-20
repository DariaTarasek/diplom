package model

import "time"

type Patient struct {
	ID          UserID    `db:"user_id"`
	FirstName   string    `db:"first_name"`
	SecondName  string    `db:"second_name"`
	Surname     *string   `db:"surname"`
	Email       *string   `db:"email"`
	BirthDate   time.Time `db:"birth_date"`
	PhoneNumber *string   `db:"phone_number"`
	Gender      string    `db:"gender"`
}
