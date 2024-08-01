package transaction

import (
	"context"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetRound(t *testing.T) {
	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	round := models.Round{
		UserID:   123,
		Bet:      models.Transaction{},
		Win:      nil,
		Finished: false,
		Refunded: false,
	}

	tests := []struct {
		name    string
		roundID models.RoundID
		before  func(uw *InMemoryRepository)
		ctx     context.Context
		expect  error
	}{
		{
			name:    "проверка существующего раунда",
			roundID: models.RoundID(roundID),
			before: func(uw *InMemoryRepository) {
				uw.transactions[models.RoundID(roundID)] = round
			},
			ctx:    nil,
			expect: nil,
		},
		{
			name:    "проверка не существующего раунда",
			roundID: models.RoundID(roundID),
			before:  func(uw *InMemoryRepository) {},
			ctx:     nil,
			expect:  ErrRoundNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewInMemoryRepository()
			tc.before(uw)
			_, err := uw.GetRound(tc.ctx, tc.roundID)
			assert.ErrorIs(t, err, tc.expect)
		})
	}
}

func TestSetWin(t *testing.T) {
	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174001")
	require.NoError(t, err)

	noTransaction, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174006")
	require.NoError(t, err)
	betID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")
	require.NoError(t, err)

	winID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174003")
	require.NoError(t, err)

	bet := &models.Transaction{
		Amount:        -10,
		TransactionID: models.TransactionID(betID),
		Created:       time.Now(),
	}

	win := models.Transaction{
		Amount:        10,
		TransactionID: models.TransactionID(winID),
		Created:       time.Now(),
	}

	tests := []struct {
		name      string
		roundID   models.RoundID
		winRound  models.Transaction
		before    func(uw *InMemoryRepository)
		ctx       context.Context
		expectErr error
	}{
		{
			name:     "проверка на добавление выигрыша",
			roundID:  models.RoundID(roundID),
			winRound: win,
			before: func(uw *InMemoryRepository) {
				uw.transactions[models.RoundID(roundID)] = models.Round{
					UserID:   123,
					Bet:      *bet,
					Win:      nil,
					Finished: false,
					Refunded: false,
				}
			},
			ctx:       nil,
			expectErr: nil,
		},
		{
			name:     "проверка на добавление в раунд с выигрышем",
			roundID:  models.RoundID(roundID),
			winRound: win,
			before: func(uw *InMemoryRepository) {
				uw.transactions[models.RoundID(roundID)] = models.Round{
					UserID:   123,
					Bet:      *bet,
					Win:      &win,
					Finished: false,
					Refunded: false,
				}
			},
			ctx:       nil,
			expectErr: ErrTransactionAlreadyExists,
		},
		{
			name:     "проверка на добавление в закрытый раунд",
			roundID:  models.RoundID(roundID),
			winRound: win,
			before: func(uw *InMemoryRepository) {
				uw.transactions[models.RoundID(roundID)] = models.Round{
					UserID:   123,
					Bet:      *bet,
					Win:      nil,
					Finished: true,
					Refunded: false,
				}
			},
			ctx:       nil,
			expectErr: ErrRoundFinished,
		},
		{
			name:    "проверка на добавление в не существующий раунд",
			roundID: models.RoundID(noTransaction),
			before: func(uw *InMemoryRepository) {
				uw.transactions[models.RoundID(roundID)] = models.Round{
					UserID:   123,
					Bet:      *bet,
					Win:      nil,
					Finished: false,
					Refunded: false,
				}
			},
			ctx:       nil,
			expectErr: ErrRoundNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewInMemoryRepository()
			tc.before(uw)
			err := uw.SetWin(tc.ctx, tc.roundID, tc.winRound)
			assert.ErrorIs(t, err, tc.expectErr)

		})
	}
}

func TestCreateBet(t *testing.T) {
	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	betID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")
	require.NoError(t, err)

	bet := &models.Transaction{
		Amount:        -10,
		TransactionID: models.TransactionID(betID),
		Created:       time.Now(),
	}

	roundBet := models.Round{
		UserID:   123,
		Bet:      *bet,
		Win:      nil,
		Finished: false,
		Refunded: false,
	}

	tests := []struct {
		name    string
		roundID models.RoundID
		round   models.Round
		before  func(uw *InMemoryRepository)
		ctx     context.Context
		expect  error
	}{
		{
			name:    "проверка существующего раунда",
			roundID: models.RoundID(roundID),
			round:   roundBet,
			before:  func(uw *InMemoryRepository) {},
			ctx:     nil,
			expect:  nil,
		},
		{
			name:    "проверка не существующего раунда",
			roundID: models.RoundID(roundID),
			before: func(uw *InMemoryRepository) {
				uw.transactions[models.RoundID(roundID)] = models.Round{
					UserID:   123,
					Bet:      *bet,
					Win:      nil,
					Finished: false,
					Refunded: false,
				}
			},
			ctx:    nil,
			expect: ErrRoundIdAlreadyExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			uw := NewInMemoryRepository()
			tc.before(uw)
			err := uw.CreateBet(tc.ctx, tc.roundID, tc.round)
			assert.ErrorIs(t, err, tc.expect)
		})
	}
}
