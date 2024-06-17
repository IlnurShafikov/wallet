package models

type UserID int

type Balance int

type Amount int

type User struct {
	ID       UserID
	Login    string
	Password []byte
}

type Wallet struct {
	UserID  UserID  `json:"user_id"`
	Balance Balance `json:"balance"`
}
