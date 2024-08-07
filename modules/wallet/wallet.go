package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/modules/wallet/request"
	"github.com/IlnurShafikov/wallet/services/transaction"
	"github.com/rs/zerolog"
	"time"
)

type Repository interface {
	Create(context.Context, models.UserID, models.Balance) error
	Get(context.Context, models.UserID) (models.Balance, error)
	Update(context.Context, models.UserID, models.Amount) (models.Balance, error)
}

var (
	ErrRefundAlreadyExists = errors.New("refund already exists")
	ErrNotRefund           = errors.New("no way to roll back transactions")
	ErrRoundIDAlready      = errors.New("round id already exist")
	ErrWinAlreadyExists    = errors.New("win already exists")
	ErrRoundFinished       = errors.New("round finished")
	ErrUpdateRoundFailed   = errors.New("update round failed")
)

type Service struct {
	walletRepository Repository
	trRepository     transaction.Repository
	log              *zerolog.Logger
}

func NewWallet(
	walletRepository Repository,
	trRepository transaction.Repository,
	logger *zerolog.Logger,
) *Service {
	return &Service{
		walletRepository: walletRepository,
		trRepository:     trRepository,
		log:              logger,
	}
}

func (w *Service) Get(
	ctx context.Context,
	userID models.UserID,
) (models.Balance, error) {
	return w.walletRepository.Get(ctx, userID)
}

func (w *Service) Update(
	ctx context.Context,
	userID models.UserID,
	amount models.Amount,
) (models.Balance, error) {
	return w.walletRepository.Update(ctx, userID, amount)
}

func (w *Service) Create(
	ctx context.Context,
	userID models.UserID,
	balance models.Balance,
) (models.Balance, error) {
	err := w.walletRepository.Create(ctx, userID, balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (w *Service) Refund(
	ctx context.Context,
	userID models.UserID,
	req request.RefundTransaction,
) (models.Balance, error) {
	round, err := w.trRepository.GetRound(ctx, req.RoundID)
	if err != nil {
		return 0, err
	}

	if round.Refunded == true {
		return 0, ErrRefundAlreadyExists
	}

	if round.Win != nil {
		return 0, ErrNotRefund
	}

	amount := round.Bet.Amount
	amount *= -1

	balance, err := w.Update(ctx, userID, amount)
	if err != nil {
		return 0, err
	}

	round.Refunded = true

	err = w.trRepository.UpdateRound(ctx, req.RoundID, *round)
	if err != nil {
		return 0, ErrUpdateRoundFailed
	}

	return balance, nil
}

func (w *Service) Change(
	ctx context.Context,
	userID models.UserID,
	req request.UpdateBalance,
) (models.Balance, error) {
	if req.IsBet() {
		return w.createBet(ctx, userID, req)
	} else if req.IsWin() {
		return w.setWin(ctx, userID, req)
	}

	return 0, errors.New("amount is not be zero")
}

func (w *Service) createBet(
	ctx context.Context,
	userID models.UserID,
	req request.UpdateBalance,
) (models.Balance, error) {
	_, err := w.trRepository.GetRound(ctx, req.RoundID)
	if err == nil {
		return 0, ErrRoundIDAlready
	}

	if !errors.Is(err, transaction.ErrRoundNotFound) {
		return 0, fmt.Errorf("get transaction: %w", err)
	}

	balance, err := w.Update(ctx, userID, req.Amount)
	if err != nil {
		return 0, fmt.Errorf("change balance: %w", err)
	}

	round := models.Round{
		UserID: userID,
		Bet: models.Transaction{
			Amount:        req.Amount,
			TransactionID: req.TransactionID,
			Created:       time.Now(),
		},
		Win:      nil,
		Finished: req.Finished,
		Refunded: false,
	}

	err = w.trRepository.CreateBet(ctx, req.RoundID, round)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (w *Service) setWin(
	ctx context.Context,
	userID models.UserID,
	req request.UpdateBalance,
) (models.Balance, error) {
	round, err := w.trRepository.GetRound(ctx, req.RoundID)
	if err != nil {
		return 0, err
	}

	if round.Finished == true {
		return 0, ErrRoundFinished
	}

	if round.Win != nil {
		return 0, ErrWinAlreadyExists
	}

	balance, err := w.Update(ctx, userID, req.Amount)
	if err != nil {
		return 0, fmt.Errorf("change balance: %w", err)
	}

	winRound := models.Transaction{
		Amount:        req.Amount,
		TransactionID: req.TransactionID,
		Created:       time.Now(),
	}

	err = w.trRepository.SetWin(ctx, req.RoundID, winRound)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
