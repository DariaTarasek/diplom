package model

type Employee struct {
	ID          int     `json:"user_id"`
	FirstName   string  `json:"firstName"`
	SecondName  string  `json:"secondName"`
	Surname     *string `json:"surname"`
	PhoneNumber *string `json:"phone"`
	Email       string  `json:"email"`
	Education   *string `json:"education"`
	Experience  *int    `json:"experience"`
	Gender      string  `json:"gender"`
	Role        int     `json:"role"`
}
