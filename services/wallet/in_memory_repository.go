package wallet

import (
	"errors"
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
)

var (
	ErrWalletNotFound           = errors.New("wallet not found")
	ErrWalletNotEnoughMoney     = errors.New("not enough money on the balance")
	ErrWalletAlreadyExists      = errors.New("wallet already exists")
	ErrWalletNotNegativeBalance = errors.New("the balance cannot be negative")
)

type InMemoryRepository struct {
	wallet map[models.UserID]models.Balance
}

// NewInMemoryRepository - создание нового экземпляра кошелька в оп
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		wallet: make(map[models.UserID]models.Balance),
	}
}

// Get - Возвращает информацию из кошелька
func (i *InMemoryRepository) Get(userID models.UserID) (models.Balance, error) {
	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}

	return balance, nil
}

// Create -  создает кошелек
func (i *InMemoryRepository) Create(userID models.UserID, balance models.Balance) error {
	if balance < 0 {
		return ErrWalletNotNegativeBalance
	}

	_, exists := i.wallet[userID]
	if exists {
		return fmt.Errorf("%w", ErrWalletAlreadyExists)
	}

	i.wallet[userID] = balance

	return nil

}

// Change - Манипуляции с балансом
func (i *InMemoryRepository) Change(userID models.UserID, amount models.Balance) (models.Balance, error) {
	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}

	balance += amount
	if balance < 0 {
		return 0, ErrWalletNotEnoughMoney
	}

	i.wallet[userID] = balance

	return balance, nil
}
