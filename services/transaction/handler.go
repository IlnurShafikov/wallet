package transaction

import (
	"context"
	"encoding/json"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/transaction/request"
	"github.com/IlnurShafikov/wallet/services/wallet/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type repository interface {
	GetRound(_ context.Context, roundID models.RoundID) (*models.Round, error)
}

type Handler struct {
	transactions repository
	log          *zerolog.Logger
}

func RunTransactionHandler(router fiber.Router, transaction repository, logger *zerolog.Logger) {
	h := &Handler{
		transactions: transaction,
		log:          logger,
	}

	transactionGroup := router.Group("/transactions")
	transactionGroup.Get("/", h.getTransaction)
}

func (h *Handler) getTransaction(fCtx *fiber.Ctx) error {
	req := request.GetTransactionRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).Msg("unmarshal failed")
		return err
	}

	round, err := h.transactions.GetRound(fCtx.Context(), req.RoundID)
	if err != nil {
		h.log.Err(err).Msg("round not found")
		return ErrRoundNotFound
	}

	var trs models.Transaction
	if req.TransactionID == round.Bet.TransactionID {
		trs = round.Bet
	} else if round.Win != nil && req.TransactionID == round.Win.TransactionID {
		trs = *round.Win
	} else {
		return ErrTransactionNotFound
	}

	h.log.Debug().
		Str("roundID", req.RoundID.String()).
		Str("transaction_id", req.TransactionID.String()).
		Msg("get transaction successful")

	err = fCtx.Status(fiber.StatusOK).JSON(response.TransactionResponse{
		Transaction: trs,
	})
	if err != nil {
		h.log.Err(err).Msg("send JSON failed")
		return err
	}

	return nil
}
