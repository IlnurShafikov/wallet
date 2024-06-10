package request

import "github.com/IlnurShafikov/wallet/models"

type CreateWallet struct {
	Balance models.Balance `json:"balance"`
}

type UpdateBalance struct {
	Amount int `json:"amount"`
}
