package model

type UserID int

type User struct {
	ID       UserID
	Login    *string
	Password *string
}
