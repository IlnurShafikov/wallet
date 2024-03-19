package response

type CreateWalletResponse struct {
	UserID  string `json:"user_id"`
	Balance int    `json:"balance"`
}

type GetBalanceResponse struct {
	Balance int `json:"balance"`
}

type UpdateBalanceResponse struct {
	UserID  string `json:"user_id"`
	Balance int    `json:"balance"`
}
