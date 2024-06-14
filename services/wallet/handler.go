package wallet

import (
	"encoding/json"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/wallet/request"
	"github.com/IlnurShafikov/wallet/services/wallet/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	wallet *InMemoryRepository
	log    *zerolog.Logger
}

func NewHandler(
	router fiber.Router,
	wallet *InMemoryRepository,
	logger *zerolog.Logger,
) *Handler {
	h := &Handler{
		wallet: wallet,
		log:    logger,
	}

	group := router.Group("/wallet")

	group.Post("/:userID", h.CreateWallet)
	group.Get("/:userID", h.GetWallet)
	group.Put("/:userID", h.UpdateBalance)

	return h
}

func (h *Handler) CreateWallet(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.CreateWallet{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).Msg("invalid request format")
		return err
	}

	err = h.wallet.Create(userID, req.Balance)
	if err != nil {
		h.log.Err(err).
			Int("userID", int(userID)).
			Msg("wallet already exists")
		return err
	}

	h.log.Info().
		Int("userID", int(userID)).
		Msg("wallet created")

	err = h.sendJson(fCtx, req.Balance, fiber.StatusCreated)
	if err != nil {
		h.log.Err(err).Msg("invalid response format")
		return err
	}

	return nil
}

func (h *Handler) GetWallet(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	balance, err := h.wallet.Get(userID)
	if err != nil {
		h.log.Err(err).Msg("wallet not found")
		return err
	}

	h.log.Info().
		Int("userID", int(userID)).
		Msg("get wallet successful")

	err = h.sendJson(fCtx, balance, fiber.StatusOK)
	if err != nil {
		return err
	}

	return nil

}

func (h *Handler) UpdateBalance(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.UpdateBalance{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).
			Int("userID", int(userID)).
			Msg("invalid request format")
		return err
	}

	balance, err := h.wallet.Change(userID, req.Amount)
	if err != nil {
		h.log.Err(err).Msg("wallet not found")
		return err
	}

	h.log.Info().
		Int("userID", int(userID)).
		Msg("change balance successful")

	err = h.sendJson(fCtx, balance, fiber.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) getUserID(fCtx *fiber.Ctx) (models.UserID, error) {
	id, err := fCtx.ParamsInt("userID")
	if err != nil {
		h.log.Err(err).Msg("error read body")
		return 0, err
	}

	return models.UserID(id), nil
}

func (h *Handler) sendJson(fCtx *fiber.Ctx, balance models.Balance, status int) error {
	err := fCtx.Status(status).JSON(response.BalanceResponse{
		Balance: balance,
	})

	if err != nil {
		h.log.Err(err).Msg("send JSON failed")
		return err
	}

	return nil
}
