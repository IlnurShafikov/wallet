package transaction

import (
	"context"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"sync"
)

var (
	ErrRoundNotFound            = errors.New("round not found")
	ErrRoundIdAlreadyExists     = errors.New("round_id already exists")
	ErrTransactionNotFound      = errors.New("transaction_id not found")
	ErrTransactionAlreadyExists = errors.New("transaction_id already exists")
	ErrRoundFinished            = errors.New("round is finished")
	ErrRoundRefundAlreadyExists = errors.New("round is refunded")
)

type InMemoryRepository struct {
	mu           sync.Mutex
	transactions map[models.RoundID]models.Round
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		transactions: make(map[models.RoundID]models.Round),
	}
}

func (i *InMemoryRepository) GetRound(_ context.Context, roundID models.RoundID) (*models.Round, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	round, exists := i.transactions[roundID]
	if !exists {
		return nil, ErrRoundNotFound
	}

	return &round, nil
}

func (i *InMemoryRepository) CreateBet(_ context.Context, roundID models.RoundID, round models.Round) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	_, exists := i.transactions[roundID]
	if exists {
		return ErrRoundIdAlreadyExists
	}

	i.transactions[roundID] = round

	return nil
}

func (i *InMemoryRepository) SetWin(_ context.Context, roundID models.RoundID, winRound models.Transaction) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	round, exists := i.transactions[roundID]
	if !exists {
		return ErrRoundNotFound
	}

	if round.Win != nil {
		return ErrTransactionAlreadyExists
	}

	if round.Finished == true {
		return ErrRoundFinished
	}

	if round.Refunded == true {
		return ErrRoundRefundAlreadyExists
	}

	i.transactions[roundID] = models.Round{
		UserID: round.UserID,
		Bet:    round.Bet,
		Win: &models.Transaction{
			Amount:        winRound.Amount,
			TransactionID: winRound.TransactionID,
			Created:       winRound.Created,
		},
		Finished: true,
		Refunded: false,
	}

	return nil
}

func (i *InMemoryRepository) UpdateRound(_ context.Context, roundID models.RoundID, updateRound models.Round) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	_, exists := i.transactions[roundID]
	if !exists {
		return ErrRoundNotFound
	}

	i.transactions[roundID] = updateRound

	return nil

}
