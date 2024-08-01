package models

import (
	"github.com/gofrs/uuid"
	"strconv"
	"time"
)

type RoundID = uuid.UUID

type TransactionID = uuid.UUID

type UserID int

func (u UserID) String() string {
	return strconv.Itoa(int(u))
}

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

// Transaction - это данные по ставке игрока
// либо по выигрышу
type Transaction struct {
	Amount        Amount        `json:"amount"`
	TransactionID TransactionID `json:"transaction_id"`
	// Время когда был создан запрос
	Created time.Time `json:"created"`
}

type Round struct {
	UserID   UserID       `json:"user_id"`
	Bet      Transaction  `json:"bet"`
	Win      *Transaction `json:"win,omitempty"`
	Finished bool         `json:"finished"`
	Refunded bool         `json:"refunded"`
}
