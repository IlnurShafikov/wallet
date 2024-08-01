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
	wallet *Wallet
	log    *zerolog.Logger
}

func RunWalletHandler(
	router fiber.Router,
	wallet *Wallet,
	logger *zerolog.Logger,
) {
	h := &Handler{
		wallet: wallet,
		log:    logger,
	}

	walletGroup := router.Group("/wallet")
	walletGroup.Post("/:userID", h.createWallet)
	walletGroup.Get("/:userID", h.getBalance)
	walletGroup.Put("/:userID", h.changeBalance)
	walletGroup.Post("refund/:userID", h.refundTransaction)
}

func (h *Handler) createWallet(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.CreateWallet{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).Msg("unmarshal failed")
		return err
	}

	_, err = h.wallet.Create(fCtx.Context(), userID, req.Balance)
	if err != nil {
		h.log.Err(err).
			Int("userID", int(userID)).
			Msg("wallet already exists")
		return err
	}

	h.log.Debug().
		Int("userID", int(userID)).
		Msg("wallet created")

	err = h.sendJson(fCtx, req.Balance, fiber.StatusCreated)
	if err != nil {
		h.log.Err(err).Msg("invalid response format")
		return err
	}

	return nil
}

func (h *Handler) getBalance(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	balance, err := h.wallet.Get(fCtx.Context(), userID)
	if err != nil {
		h.log.Err(err).Msg("wallet not found")
		return err
	}

	h.log.Debug().
		Int("userID", int(userID)).
		Msg("get wallet successful")

	err = h.sendJson(fCtx, balance, fiber.StatusOK)
	if err != nil {
		return err
	}

	return nil

}

func (h *Handler) changeBalance(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.UpdateBalance{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).
			Int("userID", int(userID)).
			Msg("unmarshal failed")
		return err
	}

	balance, err := h.wallet.Change(fCtx.Context(), userID, req)
	if err != nil {
		h.log.Err(err).Msg("wallet not found")
		return err
	}

	h.log.Debug().
		Int("userID", int(userID)).
		Str("transaction_id", req.TransactionID.String()).
		Msg("change balance successful")

	err = h.sendJson(fCtx, balance, fiber.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) refundTransaction(fCtx *fiber.Ctx) error {
	userID, err := h.getUserID(fCtx)
	if err != nil {
		return err
	}

	req := request.RefundTransaction{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).
			Int("userID", int(userID)).
			Msg("unmarshal failed")
		return err
	}

	balance, err := h.wallet.Refund(fCtx.Context(), userID, req)
	if err != nil {
		h.log.Err(err).
			Int("userID", int(userID)).
			Msg("refund failed")
		return err
	}

	h.log.Debug().
		Int("userID", int(userID)).
		Str("transaction_id", req.RoundID.String()).
		Msg("refund successful")

	err = h.sendJson(fCtx, balance, fiber.StatusOK)
	if err != nil {
		return err
	}

	return nil

}

func (h *Handler) getUserID(fCtx *fiber.Ctx) (models.UserID, error) {
	id, err := fCtx.ParamsInt("userID")
	if err != nil {
		h.log.Err(err).Msg("invalid variable type")
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
