package model

type Doctor struct {
	ID          UserID
	FirstName   string
	SecondName  string
	Surname     *string
	PhoneNumber *string
	Email       string
	Education   *string
	Experience  *int
	Gender      string
}
