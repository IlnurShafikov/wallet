package transaction

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

func NewRedisRepository(client *redis.Client, expiredAt time.Duration) *RedisRepository {
	return &RedisRepository{
		client:   client,
		expireAt: expiredAt,
	}
}

func (r *RedisRepository) UpdateRound(ctx context.Context, roundID models.RoundID, round models.Round) error {
	data, err := json.Marshal(round)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = r.client.Set(ctx, roundID.String(), data, r.expireAt).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

// Create
// Update
// Get
func (r *RedisRepository) GetRound(ctx context.Context, roundID models.RoundID) (*models.Round, error) {
	res, err := r.client.Get(ctx, roundID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = ErrRoundNotFound
		}

		return nil, err
	}

	round := new(models.Round)

	err = json.Unmarshal([]byte(res), round)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return round, nil
}

func (r *RedisRepository) CreateBet(ctx context.Context, roundID models.RoundID, round models.Round) error {
	count, err := r.client.Exists(ctx, roundID.String()).Result()
	if err != nil {
		return fmt.Errorf("redis.Exists: %w", err)
	}

	if count > 0 {
		return ErrRoundIdAlreadyExists
	}

	err = r.UpdateRound(ctx, roundID, round)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) SetWin(ctx context.Context, roundID models.RoundID, winTransaction models.Transaction) error {
	res, err := r.client.Get(ctx, roundID.String()).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = ErrRoundNotFound
		}
		return err
	}
	round := new(models.Round)

	err = json.Unmarshal([]byte(res), round)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if round.Win != nil {
		return ErrTransactionAlreadyExists
	}

	if round.Refunded == true {
		return ErrRoundRefundAlreadyExists
	}

	if round.Finished == true {
		return ErrRoundFinished
	}

	round.Win = &winTransaction
	round.Finished = true

	err = r.UpdateRound(ctx, roundID, *round)
	if err != nil {
		return err
	}

	return nil
}
