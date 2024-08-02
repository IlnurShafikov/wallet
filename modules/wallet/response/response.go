package response

import (
	"github.com/IlnurShafikov/wallet/models"
)

type BalanceResponse struct {
	Balance models.Balance `json:"balance"`
}
type TransactionResponse struct {
	Transaction models.Transaction `json:"transaction_id"`
}
