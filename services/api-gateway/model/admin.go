package model

type AdminWithRole struct {
	ID          int    `json:"user_id"`
	FirstName   string `json:"firstName"`
	SecondName  string `json:"secondName"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Role        string `json:"role"`
}

type AdminForAdminList struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Role        string `json:"role"`
}
