package model

type Doctor struct {
	ID          UserID `db:"user_id"`
	FirstName   string `db:"first_name"`
	SecondName  string `db:"second_name"`
	Surname     string `db:"surname"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	Education   string `db:"education"`
	Experience  int    `db:"experience"`
	Gender      string `db:"gender"`
}
