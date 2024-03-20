package wallet

import (
	"encoding/json"
	"github.com/IlnurShafikov/wallet/services/wallet/request"
	"github.com/IlnurShafikov/wallet/services/wallet/response"
	"io"
	"net/http"
)

type Handler struct {
	wallet *InMemory
}

func NewHandler(wallet *InMemory) *Handler {
	return &Handler{
		wallet: wallet,
	}
}

func (h *Handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error read body", http.StatusInternalServerError)
		return
	}

	req := request.CreateWalletRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := h.createWallet(req.UserID, req.Balance)
	if err != nil {
		http.Error(w, "Create wallet failed", http.StatusInternalServerError)
		return
	}

	if result.created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	resp, err := json.Marshal(&response.CreateWalletResponse{
		Balance: result.balance,
	})

	if err != nil {
		http.Error(w, "Marshal response failed", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Send response to client failed", http.StatusInternalServerError)
		return
	}
}

type createWalletResponse struct {
	balance int
	created bool
}

func (h *Handler) createWallet(userID string, balance int) (*createWalletResponse, error) {
	created := h.wallet.Create(userID)

	var err error

	var walletBalance int
	if created {
		walletBalance, err = h.wallet.Add(userID, balance)
	} else {
		walletBalance, err = h.wallet.Get(userID)
	}

	if err != nil {
		return nil, err
	}

	resp := &createWalletResponse{
		balance: walletBalance,
		created: created,
	}

	return resp, nil
}

func (h *Handler) GetWallet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error read body", http.StatusInternalServerError)
		return
	}

	req := request.GetBalanceRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	balance, err := h.wallet.Get(req.UserID)
	if err != nil {
		http.Error(w, "Error getting balance "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := response.BalanceResponse{
		Balance: balance,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}
}

func (h *Handler) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error read body", http.StatusInternalServerError)
		return
	}

	req := request.UpdateBalanceRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Error form", http.StatusInternalServerError)
		return
	}

	newBalance, err := h.wallet.Add(req.UserID, req.Amount)
	if err != nil {
		http.Error(w, "Error balance", http.StatusInternalServerError)
		return
	}

	res := response.BalanceResponse{
		Balance: newBalance,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Error to send", http.StatusInternalServerError)
		return
	}

}
