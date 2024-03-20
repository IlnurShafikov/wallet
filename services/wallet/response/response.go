package response

type CreateWalletResponse struct {
	UserID  string `json:"user_id"`
	Balance int    `json:"balance"`
}

type BalanceResponse struct {
	Balance int `json:"balance"`
}
