package request

type CreateWalletRequest struct {
	UserID  string `json:"user_id"`
	Balance int    `json:"balance"`
}

type GetBalanceRequest struct {
	UserID string `json:"user_id"`
}

type UpdateBalanceRequest struct {
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
}
