package model

type Doctor struct {
	ID          int     `json:"user_id"`
	FirstName   string  `json:"firstName"`
	SecondName  string  `json:"secondName"`
	Surname     *string `json:"surname"`
	PhoneNumber *string `json:"phone"`
	Email       string  `json:"email"`
	Education   *string `json:"education"`
	Experience  *int    `json:"experience"`
	Gender      string  `json:"gender"`
}

type DoctorWithSpecs struct {
	ID          int    `json:"user_id"`
	FirstName   string `json:"firstName"`
	SecondName  string `json:"secondName"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone"`
	Email       string `json:"email"`
	Education   string `json:"education"`
	Experience  int    `json:"experience"`
	Gender      string `json:"gender"`
	Specs       []int  `json:"specialty"`
}

type Spec struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
