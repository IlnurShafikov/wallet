package wallet

import (
	"context"
	"encoding/json"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRedisRepository_Get(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	userID := models.UserID(1992)
	balance := models.Balance(100)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	expErr := func(expErr error) func(t *testing.T, res models.Balance, err error) {
		return func(t *testing.T, res models.Balance, err error) {
			assert.ErrorIs(t, expErr, err)
			assert.Zero(t, res)
		}
	}

	tests := []struct {
		name     string
		userID   models.UserID
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, res models.Balance, err error)
	}{
		{
			name:     "get wallet: wallet not found",
			userID:   userID,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(ErrWalletNotFound),
		}, {
			name:   "get wallet successfully",
			userID: userID,
			before: func(t *testing.T, r *redis.Client) {
				balanceJSON, err := json.Marshal(balance)
				require.NoError(t, err)

				err = r.Set(ctx, userID.String(), balanceJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, res models.Balance, err error) {
				require.Equal(t, res, balance)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			balance, err := repo.Get(ctx, tc.userID)
			tc.checkRes(t, balance, err)
		})
	}
}

func TestRedisRepository_Create(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	userID := models.UserID(1992)
	balance := models.Balance(100)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	repo := &RedisRepository{
		client: client,
	}

	expErr := func(expErr error) func(t *testing.T, err error) {
		return func(t *testing.T, err error) {
			assert.ErrorIs(t, expErr, err)
		}
	}

	tests := []struct {
		name     string
		userID   models.UserID
		balance  models.Balance
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, err error)
	}{
		{
			name:     "create wallet: balance not negative",
			userID:   userID,
			balance:  -balance,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(ErrWalletNotNegativeBalance),
		}, {
			name:    "create wallet: wallet already exists",
			userID:  userID,
			balance: balance,
			before: func(t *testing.T, r *redis.Client) {
				balanceJSON, err := json.Marshal(balance)
				require.NoError(t, err)

				err = r.Set(ctx, userID.String(), balanceJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrWalletAlreadyExists),
		}, {
			name:     "create wallet successfully",
			userID:   userID,
			balance:  balance,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			err = repo.Create(ctx, tc.userID, tc.balance)
			tc.checkRes(t, err)
		})
	}
}

func TestRedisRepository_Update(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	userID := models.UserID(1992)
	balance := models.Balance(100)
	amount := models.Amount(10)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	expErr := func(expErr error) func(t *testing.T, res models.Balance, err error) {
		return func(t *testing.T, res models.Balance, err error) {
			assert.ErrorIs(t, expErr, err)
			assert.Zero(t, res)
		}
	}

	tests := []struct {
		name     string
		userID   models.UserID
		amount   models.Amount
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, res models.Balance, err error)
	}{
		{
			name:     "update wallet: wallet not found",
			userID:   userID,
			amount:   amount,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(ErrWalletNotFound),
		}, {
			name:   "update wallet: wallet not enough money",
			userID: userID,
			amount: models.Amount(-200),
			before: func(t *testing.T, r *redis.Client) {
				balanceJSON, err := json.Marshal(balance)
				require.NoError(t, err)

				err = r.Set(ctx, userID.String(), balanceJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrWalletNotEnoughMoney),
		}, {
			name:   "update wallet: positive amount successfully",
			userID: userID,
			amount: amount,
			before: func(t *testing.T, r *redis.Client) {
				balanceJSON, err := json.Marshal(balance)
				require.NoError(t, err)

				err = r.Set(ctx, userID.String(), balanceJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, res models.Balance, err error) {
				require.Equal(t, res, models.Balance(110))
			},
		}, {
			name:   "update wallet: negative amount successfully",
			userID: userID,
			amount: -amount,
			before: func(t *testing.T, r *redis.Client) {
				balanceJSON, err := json.Marshal(balance)
				require.NoError(t, err)

				err = r.Set(ctx, userID.String(), balanceJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, res models.Balance, err error) {
				require.Equal(t, res, models.Balance(90))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			balance, err := repo.Update(ctx, tc.userID, tc.amount)
			tc.checkRes(t, balance, err)
		})
	}
}
