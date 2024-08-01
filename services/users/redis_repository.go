package users

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/go-redis/redis/v8"
	"time"
)

var _ = Repository(&RedisRepository{})

type RedisRepository struct {
	client   *redis.Client
	expireAt time.Duration
	lastID   models.UserID
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client:   client,
		expireAt: 0,
		lastID:   0,
	}
}

func (r *RedisRepository) Create(_ context.Context, login string, password []byte) (*models.User, error) {
	ctx := context.Background()
	count, err := r.client.Exists(ctx, login).Result()
	if err != nil {
		return nil, fmt.Errorf("redis.Exists: %w", err)
	}

	if count > 0 {
		return nil, ErrUserAlreadyExists
	}
	r.lastID++

	user := models.User{
		ID:       r.lastID,
		Login:    login,
		Password: password,
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	err = r.client.Set(ctx, login, data, r.expireAt).Err()

	return &user, nil
}

func (r *RedisRepository) Get(_ context.Context, login string) (*models.User, error) {
	ctx := context.Background()
	res, err := r.client.Get(ctx, login).Result()
	if err != nil {
		return nil, ErrUserNotFound
	}

	user := new(models.User)

	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return user, nil
}
