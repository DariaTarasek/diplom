package model

type Doctor struct {
	ID          int     `json:"user_id"`
	FirstName   string  `json:"first_name"`
	SecondName  string  `json:"second_name"`
	Surname     *string `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
	Email       string  `json:"email"`
	Education   *string `json:"education"`
	Experience  *int    `json:"experience"`
	Gender      string  `json:"gender"`
}
