package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository struct {
	client   *redis.Client
	expireAt time.Duration
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client:   client,
		expireAt: 0,
	}
}

func (r *RedisRepository) Get(
	ctx context.Context,
	userID models.UserID,
) (models.Balance, error) {
	res, err := r.client.Get(ctx, userID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = ErrWalletNotFound
		}
		return 0, err
	}

	balance := new(models.Balance)

	err = json.Unmarshal([]byte(res), &balance)
	if err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	return *balance, nil
}

func (r *RedisRepository) Create(
	ctx context.Context,
	userID models.UserID,
	balance models.Balance,
) error {
	if balance < 0 {
		return ErrWalletNotNegativeBalance
	}

	count, err := r.client.Exists(ctx, userID.String()).Result()
	if err != nil {
		return fmt.Errorf("redis.Exists: %w", err)
	}

	if count > 0 {
		return ErrWalletAlreadyExists
	}

	data, err := json.Marshal(balance)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = r.client.Set(ctx, userID.String(), data, r.expireAt).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (r *RedisRepository) Update(
	ctx context.Context,
	userID models.UserID,
	amount models.Amount,
) (models.Balance, error) {
	res, err := r.client.Get(ctx, userID.String()).Result()
	if err != nil {
		return 0, ErrWalletNotFound
	}

	balance := new(models.Balance)

	err = json.Unmarshal([]byte(res), balance)
	if err != nil {
		return 0, fmt.Errorf("unmarshal: %w", err)
	}

	*balance += models.Balance(amount)

	if *balance < 0 {
		return 0, ErrWalletNotEnoughMoney
	}

	data, err := json.Marshal(balance)
	if err != nil {
		return 0, fmt.Errorf("marshal: %w", err)
	}

	err = r.client.Set(ctx, userID.String(), data, r.expireAt).Err()
	if err != nil {
		return 0, fmt.Errorf("redis.Set: %w", err)
	}

	return *balance, nil
}
