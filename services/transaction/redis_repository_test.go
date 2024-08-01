package transaction

import (
	"context"
	"encoding/json"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisRepository_GetRound(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	const (
		amount = -10
	)

	userID := models.UserID(1992)
	ctx := context.Background()

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	roundBet := &models.Round{
		UserID: userID,
		Bet: models.Transaction{
			Amount:        amount,
			TransactionID: models.RoundID(roundID),
			Created:       time.Time{},
		},
		Win:      nil,
		Finished: false,
		Refunded: false,
	}

	expErr := func(expErr error) func(t *testing.T, res *models.Round, err error) {
		return func(t *testing.T, res *models.Round, err error) {
			assert.ErrorIs(t, expErr, err)
			assert.Nil(t, res)
		}
	}

	tests := []struct {
		name     string
		roundID  models.RoundID
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, req *models.Round, err error)
	}{
		{
			name:     "get round: wallet not found",
			roundID:  models.RoundID(roundID),
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(ErrRoundNotFound),
		}, {
			name:    "get round: unmarshal error",
			roundID: models.RoundID(roundID),
			before: func(t *testing.T, r *redis.Client) {
				err := r.Set(ctx, roundID.String(), "invalid json", 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, req *models.Round, err error) {
				require.ErrorContains(t, err, "unmarshal: invalid character")
			},
		}, {
			name:    "success",
			roundID: models.RoundID(roundID),
			before: func(t *testing.T, r *redis.Client) {
				roundJSON, err := json.Marshal(roundBet)
				require.NoError(t, err)
				err = r.Set(ctx, roundID.String(), roundJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, req *models.Round, err error) {
				require.Equal(t, roundBet, req)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			round, err := repo.GetRound(ctx, tc.roundID)
			tc.checkRes(t, round, err)
		})
	}

}

func TestRedisRepository_CreateBet(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	const (
		bet = -10
	)

	userID := models.UserID(1992)
	ctx := context.Background()

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	roundBet := &models.Round{
		UserID: userID,
		Bet: models.Transaction{
			Amount:        bet,
			TransactionID: models.RoundID(roundID),
			Created:       time.Time{},
		},
		Win:      nil,
		Finished: false,
		Refunded: false,
	}

	expErr := func(exErr error) func(t *testing.T, err error) {
		return func(t *testing.T, err error) {
			assert.ErrorIs(t, exErr, err)
		}
	}

	tests := []struct {
		name     string
		roundID  models.RoundID
		round    models.Round
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, err error)
	}{
		{
			name:    "create bet: round already exists",
			roundID: models.RoundID(roundID),
			round:   *roundBet,
			before: func(t *testing.T, r *redis.Client) {
				roundJSON, err := json.Marshal(roundBet)
				require.NoError(t, err)
				err = r.Set(ctx, roundID.String(), roundJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrRoundIdAlreadyExists),
		}, {
			name:     "create bet successfully",
			roundID:  models.RoundID(roundID),
			round:    *roundBet,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			err = repo.CreateBet(ctx, tc.roundID, tc.round)
			tc.checkRes(t, err)

		})
	}
}

func TestInMemoryRepository_SetWin(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	const (
		bet = -10
		win = 12
	)

	userID := models.UserID(1992)
	ctx := context.Background()

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	newRoundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174012")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	winTransaction := models.Transaction{
		Amount:        win,
		TransactionID: models.TransactionID(roundID),
		Created:       time.Time{},
	}

	expErr := func(expErr error) func(t *testing.T, err error) {
		return func(t *testing.T, err error) {
			assert.ErrorIs(t, expErr, err)
		}
	}

	makeBet := func(finished, refunded bool) *models.Round {
		return &models.Round{
			UserID: userID,
			Bet: models.Transaction{
				Amount:        bet,
				TransactionID: models.TransactionID{},
				Created:       time.Time{},
			},
			Win:      nil,
			Finished: finished,
			Refunded: refunded,
		}
	}

	setWin := func(tr *models.Round) {
		tr.Win = &models.Transaction{
			Amount:        win,
			TransactionID: models.TransactionID{},
			Created:       time.Time{},
		}
		tr.Finished = true
	}

	test := []struct {
		name           string
		roundID        models.RoundID
		winTransaction models.Transaction
		before         func(t *testing.T, r *redis.Client)
		checkRes       func(t *testing.T, err error)
	}{
		{
			name:           "set win: round not found",
			roundID:        models.RoundID(roundID),
			winTransaction: winTransaction,
			before:         func(t *testing.T, r *redis.Client) {},
			checkRes:       expErr(ErrRoundNotFound),
		},
		{
			name:           "set win: unmarshal failed",
			roundID:        models.RoundID(roundID),
			winTransaction: winTransaction,
			before: func(t *testing.T, r *redis.Client) {
				err = r.Set(ctx, roundID.String(), "invalid json", 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "unmarshal:")
			},
		},
		{
			name:           "set win: win transaction already exists",
			roundID:        models.RoundID(roundID),
			winTransaction: winTransaction,
			before: func(t *testing.T, r *redis.Client) {
				round := makeBet(false, false)
				setWin(round)

				roundJSON, err := json.Marshal(round)
				require.NoError(t, err)

				err = r.Set(ctx, roundID.String(), roundJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrTransactionAlreadyExists),
		},
		{
			name:           "set win: round refunded",
			roundID:        models.RoundID(roundID),
			winTransaction: winTransaction,
			before: func(t *testing.T, r *redis.Client) {
				round := makeBet(false, true)

				roundJSON, err := json.Marshal(round)
				require.NoError(t, err)

				err = r.Set(ctx, roundID.String(), roundJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrRoundRefundAlreadyExists),
		},
		{
			name:           "set win: round finished",
			roundID:        models.RoundID(roundID),
			winTransaction: winTransaction,
			before: func(t *testing.T, r *redis.Client) {
				round := makeBet(true, false)
				roundJSON, err := json.Marshal(round)
				require.NoError(t, err)

				err = r.Set(ctx, roundID.String(), roundJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrRoundFinished),
		},
		{
			name:           "set win: marshal failed",
			roundID:        models.RoundID(roundID),
			winTransaction: winTransaction,
			before: func(t *testing.T, r *redis.Client) {
				err := r.Set(ctx, roundID.String(), "invalid json", 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, err error) {
				require.ErrorContains(t, err, "marshal:")
			},
		},
		{
			name:           "set win successful",
			roundID:        models.RoundID(newRoundID),
			winTransaction: winTransaction,
			before: func(t *testing.T, r *redis.Client) {
				round := makeBet(false, false)

				roundJson, err := json.Marshal(round)
				require.NoError(t, err)

				err = r.Set(ctx, newRoundID.String(), roundJson, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(nil),
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			err = repo.SetWin(ctx, tc.roundID, tc.winTransaction)
			tc.checkRes(t, err)
		})
	}
}

func TestInMemoryRepository_UpdateRound(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	const (
		amount = -10
	)

	var testTime = time.Date(2023, time.July, 26, 15, 30, 0, 0, time.UTC)

	userID := models.UserID(1992)
	ctx := context.Background()

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	roundBet := &models.Round{
		UserID: userID,
		Bet: models.Transaction{
			Amount:        amount,
			TransactionID: models.RoundID(roundID),
			Created:       testTime,
		},
		Win:      nil,
		Finished: false,
		Refunded: false,
	}

	winTransaction := models.Transaction{
		Amount:        -amount,
		TransactionID: models.TransactionID(roundID),
		Created:       testTime,
	}

	updateRound := &models.Round{
		UserID: userID,
		Bet: models.Transaction{
			Amount:        amount,
			TransactionID: models.RoundID(roundID),
			Created:       testTime,
		},
		Win:      &winTransaction,
		Finished: false,
		Refunded: false,
	}

	expErr := func(expErr error) func(t *testing.T, err error) {
		return func(t *testing.T, err error) {
			assert.ErrorIs(t, expErr, err)
		}
	}

	tests := []struct {
		name        string
		roundID     models.RoundID
		updateRound models.Round
		before      func(t *testing.T, r *redis.Client)
		checkRes    func(t *testing.T, err error)
	}{
		{
			name:        "update round: round not found",
			roundID:     models.RoundID(roundID),
			updateRound: *updateRound,
			before:      func(t *testing.T, r *redis.Client) {},
			checkRes:    expErr(ErrRoundNotFound),
		},
		{
			name:        "update round successfully",
			roundID:     models.RoundID(roundID),
			updateRound: *updateRound,
			before: func(t *testing.T, r *redis.Client) {
				roundJSON, err := json.Marshal(roundBet)
				require.NoError(t, err)

				err = r.Set(ctx, roundID.String(), roundJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			err = repo.UpdateRound(ctx, tc.roundID, tc.updateRound)
			tc.checkRes(t, err)
		})
	}
}
