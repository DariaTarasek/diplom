package model

type Doctor struct {
	ID          int
	FirstName   string
	SecondName  string
	Surname     string
	PhoneNumber string
	Email       string
	Education   string
	Experience  int
	Gender      string
	Specs       []int
}

type Spec struct {
	ID   int
	Name string
}
