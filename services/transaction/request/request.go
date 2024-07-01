package request

import (
	"github.com/IlnurShafikov/wallet/models"
)

type GetTransactionRequest struct {
	RoundID       models.RoundID       `json:"round_id"`
	TransactionID models.TransactionID `json:"transaction_id"`
}
