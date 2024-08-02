package wallet

import (
	"context"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"sync"
)

var (
	ErrWalletNotFound           = errors.New("wallet not found")
	ErrWalletNotEnoughMoney     = errors.New("not enough money on the balance")
	ErrWalletAlreadyExists      = errors.New("wallet already exists")
	ErrWalletNotNegativeBalance = errors.New("the balance cannot be negative")
)

type InMemoryRepository struct {
	mu     sync.Mutex
	wallet map[models.UserID]models.Balance
}

// NewInMemoryRepository - создание нового экземпляра кошелька в оп
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		wallet: make(map[models.UserID]models.Balance),
	}
}

// Get - Возвращает информацию из кошелька
func (i *InMemoryRepository) Get(_ context.Context, userID models.UserID) (models.Balance, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}

	return balance, nil
}

// Create -  создает кошелек
func (i *InMemoryRepository) Create(
	_ context.Context,
	userID models.UserID,
	balance models.Balance,
) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if balance < 0 {
		return ErrWalletNotNegativeBalance
	}

	_, exists := i.wallet[userID]
	if exists {
		return ErrWalletAlreadyExists
	}

	i.wallet[userID] = balance

	return nil
}

// Update - Манипуляции с балансом
func (i *InMemoryRepository) Update(
	_ context.Context,
	userID models.UserID,
	amount models.Amount,
) (models.Balance, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}

	balance += models.Balance(amount)
	if balance < 0 {
		return 0, ErrWalletNotEnoughMoney
	}

	i.wallet[userID] = balance

	return balance, nil
}
