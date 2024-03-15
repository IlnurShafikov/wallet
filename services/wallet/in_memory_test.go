package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	const userID = "1"

	tests := []struct {
		name   string
		userID string
		before func(uw *InMemory)
		expect bool
	}{
		{
			name:   "добавление пользователя которого не существует",
			userID: userID,
			expect: true,
			before: func(uw *InMemory) {},
		},
		{
			name:   "добавление пользователя который существует",
			userID: userID,
			before: func(uw *InMemory) {
				uw.wallet[userID] = 0
			},
			expect: false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			uw := NewWallet()
			testCase.before(uw)
			got := uw.Create(testCase.userID)
			assert.Equal(t, testCase.expect, got)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		expect    int
		expectErr error
	}{
		{
			name:      "Проверяем что пользователь существует",
			userID:    "user01",
			expect:    0,
			expectErr: nil,
		},
		{
			name:      "Проверяем что пользователь не существует",
			userID:    "",
			expect:    0,
			expectErr: ErrWalletNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewWallet()
			_ = uw.Create("user01")
			got, err := uw.Get(tc.userID)
			assert.Equal(t, tc.expect, got)
			assert.Equal(t, tc.expectErr, err)
		})

	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		userID    string
		balance   int
		amount    int
		expect    int
		expectErr error
	}{
		{
			userID:    "user01",
			balance:   0,
			amount:    10,
			expect:    10,
			expectErr: nil,
		},
		{
			userID:    "user01",
			balance:   10,
			amount:    10,
			expect:    20,
			expectErr: nil,
		},
		{
			userID:    "user01",
			balance:   20,
			amount:    -20,
			expect:    0,
			expectErr: nil,
		},
		{
			userID:    "user02",
			amount:    0,
			expect:    0,
			expectErr: ErrWalletNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.userID, func(t *testing.T) {
			uw := NewWallet()
			_ = uw.Create("user01")
			got, err := uw.Add(tc.userID, tc.amount+tc.balance)
			assert.Equal(t, tc.expect, got)
			assert.Equal(t, tc.expectErr, err)
		})

	}
}
