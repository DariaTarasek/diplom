package model

type UserID int

type User struct {
	ID       UserID `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}
