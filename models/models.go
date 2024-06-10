package models

type UserID int

type Balance int

type User struct {
	UserID   UserID
	Login    string
	Password []byte
}

type Wallet struct {
	UserID  UserID  `json:"user_id"`
	Balance Balance `json:"balance"`
}
