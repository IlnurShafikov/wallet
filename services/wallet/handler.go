package wallet

import (
	"encoding/json"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/wallet/request"
	"github.com/IlnurShafikov/wallet/services/wallet/response"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	wallet *InMemoryRepository
}

func NewHandler(
	router fiber.Router,
	wallet *InMemoryRepository,
) *Handler {
	h := &Handler{
		wallet: wallet,
	}

	group := router.Group("/wallet")

	group.Post("/:userID", h.CreateWallet)
	group.Get("/:userID", h.GetWallet)
	group.Put("/:userID", h.UpdateBalance)

	return h
}

func (h *Handler) CreateWallet(fCtx *fiber.Ctx) error {
	userID, err := getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.CreateWallet{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		return err
	}

	err = h.wallet.Create(userID, req.Balance)
	if err != nil {
		return err
	}

	err = fCtx.Status(fiber.StatusCreated).JSON(response.BalanceResponse{
		Balance: req.Balance,
	})

	return err
}

func (h *Handler) GetWallet(fCtx *fiber.Ctx) error {
	userID, err := getUserID(fCtx)
	if err != nil {
		return err
	}

	balance, err := h.wallet.Get(userID)
	if err != nil {
		return err
	}

	err = fCtx.Status(fiber.StatusOK).JSON(response.BalanceResponse{
		Balance: balance,
	})

	return err

}

func (h *Handler) UpdateBalance(fCtx *fiber.Ctx) error {
	userID, err := getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.UpdateBalance{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		return err
	}

	newBalance, err := h.wallet.Change(userID, req.Amount)
	if err != nil {
		return err
	}

	err = fCtx.Status(fiber.StatusOK).JSON(response.BalanceResponse{
		Balance: newBalance,
	})

	return err
}

func getUserID(fCtx *fiber.Ctx) (models.UserID, error) {
	id, err := fCtx.ParamsInt("userID")
	if err != nil {
		return 0, err
	}

	return models.UserID(id), nil
}
