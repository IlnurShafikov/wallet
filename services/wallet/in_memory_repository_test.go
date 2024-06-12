package wallet

import (
	"github.com/IlnurShafikov/wallet/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	const userID = 123

	tests := []struct {
		name    string
		balance models.Balance
		before  func(uw *InMemoryRepository)
		expect  error
	}{
		{
			name:    "создание нового кошелька",
			expect:  nil,
			balance: 0,
			before:  func(uw *InMemoryRepository) {},
		},
		{
			name:    "создание существующего кошелька",
			balance: 10,
			before: func(uw *InMemoryRepository) {
				uw.wallet[userID] = 0
			},
			expect: ErrWalletAlreadyExists,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			uw := NewInMemoryRepository()
			testCase.before(uw)
			err := uw.Create(userID, testCase.balance)
			assert.ErrorIs(t, err, testCase.expect)
		})
	}
}

func TestGet(t *testing.T) {
	const (
		userID  = 123
		balance = 10
	)
	tests := []struct {
		name      string
		expect    models.Balance
		expectErr error
		before    func(uw *InMemoryRepository)
	}{
		{
			name:      "Получаем существующий кошелек",
			expect:    balance,
			expectErr: nil,
			before: func(uw *InMemoryRepository) {
				uw.wallet[userID] = balance
			},
		},
		{
			name:      "Получаем не существующий кошелек",
			expect:    0,
			expectErr: ErrWalletNotFound,
			before: func(_ *InMemoryRepository) {

			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewInMemoryRepository()
			tc.before(uw)
			got, err := uw.Get(userID)
			assert.Equal(t, tc.expect, got)
			assert.ErrorIs(t, err, tc.expectErr)
		})

	}
}

func TestAdd(t *testing.T) {
	const (
		userID  = 123
		balance = 10
	)
	tests := []struct {
		name      string
		amount    models.Balance
		before    func(uw *InMemoryRepository)
		expect    models.Balance
		expectErr error
	}{
		{
			name:      "test1",
			amount:    10,
			expect:    10,
			expectErr: nil,
			before: func(uw *InMemoryRepository) {
				uw.wallet[userID] = 0
			},
		},
		{
			name:      "test2",
			amount:    10,
			expect:    20,
			expectErr: nil,
			before: func(uw *InMemoryRepository) {
				uw.wallet[userID] = balance
			},
		},
		{
			name:      "test3",
			amount:    -20,
			expect:    0,
			expectErr: nil,
			before: func(uw *InMemoryRepository) {
				uw.wallet[userID] = 20
			},
		},
		{
			name:      "test4",
			amount:    0,
			expect:    0,
			expectErr: ErrWalletNotFound,
			before: func(uw *InMemoryRepository) {
				uw.wallet[133] = 0
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewInMemoryRepository()
			tc.before(uw)
			got, err := uw.Change(userID, tc.amount)
			assert.Equal(t, tc.expect, got)
			assert.ErrorIs(t, err, tc.expectErr)
		})

	}
}
