package wallet

import "errors"

var ErrWalletNotFound = errors.New("wallet not found")

type InMemory struct {
	wallet map[string]int
}
//создание нового экземпляра кошелька в оп
func NewWallet() *InMemory {
	return &InMemory{
		wallet: make(map[string]int),
	}
}
// Возвращает информацию из кошелька
func(i *InMemory) Get(userID string) (int, error) {
	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}

	return balance, nil
}
//  создает кошелек
 func (i *InMemory) Create(userID string) bool {
	_, exists := i.wallet[userID] 
	if exists {
		return false
	}
	i.wallet[userID] = 0
	return true

 }

// Манипуляции с балансом

func(i *InMemory) Add(userID string, amount int) (int, error) {
	balance, ok := i.wallet[userID]
	if !ok {
		return 0, ErrWalletNotFound
	}
	balance += amount
	i.wallet[userID] = balance
	return balance, nil
}

