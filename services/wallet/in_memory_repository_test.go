package wallet

import (
	"github.com/IlnurShafikov/wallet/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	const userID models.UserID = 1

	tests := []struct {
		name    string
		userID  models.UserID
		balance models.Balance
		before  func(uw *InMemoryRepository)
		expect  bool
	}{
		{
			name:    "добавление пользователя которого не существует",
			userID:  userID,
			expect:  true,
			balance: 0,
			before:  func(uw *InMemoryRepository) {},
		},
		{
			name:    "добавление пользователя который существует",
			userID:  userID,
			balance: 10,
			before: func(uw *InMemoryRepository) {
				uw.wallet[userID] = 0
			},
			expect: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			uw := NewWalletRepository()
			testCase.before(uw)
			got := uw.Create(testCase.userID, testCase.balance)
			assert.Equal(t, testCase.expect, got)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		userID    models.UserID
		expect    int
		expectErr error
	}{
		{
			name:      "Проверяем что пользователь существует",
			userID:    123,
			expect:    0,
			expectErr: nil,
		},
		{
			name:      "Проверяем что пользователь не существует",
			userID:    0,
			expect:    0,
			expectErr: ErrWalletNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewWalletRepository()
			_ = uw.Create(123, 10)
			got, err := uw.Get(tc.userID)
			assert.Equal(t, tc.expect, got)
			assert.Equal(t, tc.expectErr, err)
		})

	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name      string
		userID    models.UserID
		balance   models.Balance
		amount    int
		expect    models.Balance
		expectErr error
	}{
		{
			name:      "test1",
			userID:    123,
			balance:   0,
			amount:    10,
			expect:    10,
			expectErr: nil,
		},
		{
			name:      "test2",
			userID:    123,
			balance:   10,
			amount:    10,
			expect:    20,
			expectErr: nil,
		},
		{
			name:      "test3",
			userID:    123,
			balance:   20,
			amount:    -20,
			expect:    0,
			expectErr: nil,
		},
		{
			name:      "test4",
			userID:    112,
			amount:    0,
			expect:    0,
			expectErr: ErrWalletNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewWalletRepository()
			_ = uw.Create(123, 0)
			got, err := uw.Change(tc.userID, tc.amount)
			assert.Equal(t, tc.expect, got)
			assert.Equal(t, tc.expectErr, err)
		})

	}
}
