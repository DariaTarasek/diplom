package model

type Patient struct {
	ID          int     `json:"user_id"`
	FirstName   string  `json:"firstName"`
	SecondName  string  `json:"secondName"`
	Surname     *string `json:"surname"`
	PhoneNumber string  `json:"phone"`
	Email       *string `json:"email"`
	BirthDate   string  `json:"birthDate"`
	Gender      string  `json:"gender"`
	Password    string  `json:"password"`
}

type PatientWithoutPassword struct {
	ID          int     `json:"user_id"`
	FirstName   string  `json:"firstName"`
	SecondName  string  `json:"secondName"`
	Surname     *string `json:"surname"`
	PhoneNumber *string `json:"phone"`
	Email       *string `json:"email"`
	BirthDate   string  `json:"birthDate"`
	Gender      string  `json:"gender"`
}
