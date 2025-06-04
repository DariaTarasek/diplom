package model

type (
	Admin struct {
		ID          UserID
		FirstName   string
		SecondName  string
		Surname     *string
		PhoneNumber *string
		Email       string
		Gender      string
	}
	AdminWithRole struct {
		ID          UserID
		FirstName   string
		SecondName  string
		Surname     *string
		PhoneNumber *string
		Email       string
		Gender      string
		Role        string
	}
)
