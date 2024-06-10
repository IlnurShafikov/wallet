package wallet

import (
	"errors"
	"github.com/IlnurShafikov/wallet/models"
)

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrWalletNotEnoughMoney = errors.New("not enough money on the balance")
	ErrWalletAlreadyExists  = errors.New("wallet already exists")
)

type InMemoryRepository struct {
	wallet map[models.UserID]models.Balance
}

// NewWalletRepository - создание нового экземпляра кошелька в оп
func NewWalletRepository() *InMemoryRepository {
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
	_, exists := i.wallet[userID]
	if exists {
		return ErrWalletAlreadyExists
	}

	i.wallet[userID] = balance

	return nil

}

// Change - Манипуляции с балансом
func (i *InMemoryRepository) Change(userID models.UserID, amount int) (models.Balance, error) {
	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}

	newBalance := models.Balance(int(balance) + amount)
	if newBalance < 0 {
		return 0, ErrWalletNotEnoughMoney
	}

	i.wallet[userID] = newBalance

	return newBalance, nil
}
