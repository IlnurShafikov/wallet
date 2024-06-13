package response

import "github.com/IlnurShafikov/wallet/models"

type BalanceResponse struct {
	Balance models.Balance `json:"balance"`
}
