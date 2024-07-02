package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/transaction"
	"github.com/IlnurShafikov/wallet/services/wallet/request"
	"github.com/rs/zerolog"
	"time"
)

type walletRepository interface {
	Create(context.Context, models.UserID, models.Balance) error
	Get(context.Context, models.UserID) (models.Balance, error)
	Update(context.Context, models.UserID, models.Amount) (models.Balance, error)
}

type transactionRepository interface {
	GetRound(context.Context, models.RoundID) (*models.Round, error)
	CreateBet(context.Context, models.RoundID, models.Round) error
	SetWin(context.Context, models.RoundID, models.Transaction) error
	UpdateRound(context.Context, models.RoundID, models.Round) error
}

type Wallet struct {
	walletRepository walletRepository
	trRepository     transactionRepository
	log              *zerolog.Logger
}

func NewWallet(
	walletRepository walletRepository,
	trRepository transactionRepository,
	logger *zerolog.Logger,
) *Wallet {
	return &Wallet{
		walletRepository: walletRepository,
		trRepository:     trRepository,
		log:              logger,
	}
}

func (w *Wallet) Get(
	ctx context.Context,
	userID models.UserID,
) (models.Balance, error) {
	return w.walletRepository.Get(ctx, userID)
}

func (w *Wallet) Update(
	ctx context.Context,
	userID models.UserID,
	amount models.Amount,
) (models.Balance, error) {
	return w.walletRepository.Update(ctx, userID, amount)
}

func (w *Wallet) Create(
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

func (w *Wallet) Refund(
	ctx context.Context,
	userID models.UserID,
	req request.RefundTransaction,
) (models.Balance, error) {
	round, err := w.trRepository.GetRound(ctx, req.RoundID)
	if err != nil {
		return 0, transaction.ErrRoundNotFound
	}

	if round.Refunded == true {
		return 0, transaction.ErrRefundAlreadyExists
	}

	if round.Win != nil {
		return 0, transaction.ErrNotRefund
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
		return 0, err
	}

	fmt.Println(round.Refunded, round.Finished)

	return balance, nil
}

func (w *Wallet) Change(
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

func (w *Wallet) createBet(
	ctx context.Context,
	userID models.UserID,
	req request.UpdateBalance,
) (models.Balance, error) {
	_, err := w.trRepository.GetRound(ctx, req.RoundID)
	if err == nil {
		return 0, errors.New("roundID exist")
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
		return 0, fmt.Errorf("create bet transaction: %w", err)
	}

	return balance, nil
}

func (w *Wallet) setWin(
	ctx context.Context,
	userID models.UserID,
	req request.UpdateBalance,
) (models.Balance, error) {
	round, err := w.trRepository.GetRound(ctx, req.RoundID)
	if err != nil {
		return 0, err
	}

	if round.Bet.Amount == 0 {
		return 0, err
	}

	if round.Finished == true {
		return 0, transaction.ErrRoundFinished
	}

	if round.Win != nil {
		return 0, transaction.ErrTransactionAlreadyExists
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
