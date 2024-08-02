package request

import (
	"github.com/IlnurShafikov/wallet/models"
)

type CreateWallet struct {
	Balance models.Balance `json:"balance"`
}

type UpdateBalance struct {
	Amount        models.Amount        `json:"amount"`
	RoundID       models.RoundID       `json:"round_id"`
	TransactionID models.TransactionID `json:"transaction_id"`
	Finished      bool                 `json:"finished"`
}

func (u UpdateBalance) IsBet() bool {
	return u.Amount < 0
}

func (u UpdateBalance) IsWin() bool {
	return u.Amount > 0 ||
		u.Amount == 0 && u.Finished
}

type RefundTransaction struct {
	RoundID models.RoundID `json:"round_id"`
}
