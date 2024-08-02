package repositories

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

	password := []byte("1231231")
	login := "user1"
	lastID := models.UserID(12)
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	user := &models.User{
		ID:       lastID,
		Login:    login,
		Password: password,
	}

	expErr := func(expErr error) func(t *testing.T, res *models.User, err error) {
		return func(t *testing.T, res *models.User, err error) {
			assert.ErrorIs(t, expErr, err)
			assert.Nil(t, res)
		}
	}

	tests := []struct {
		name     string
		login    string
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, res *models.User, err error)
	}{
		{
			name:     "get user: user not found",
			login:    login,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: expErr(ErrUserNotFound),
		}, {
			name:  "get user successfully",
			login: login,
			before: func(t *testing.T, r *redis.Client) {
				userJSON, err := json.Marshal(user)
				require.NoError(t, err)

				err = r.Set(ctx, login, userJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: func(t *testing.T, res *models.User, err error) {
				require.Equal(t, user, res)
			},
		}, {
			name:  "get wrong user",
			login: "user2",
			before: func(t *testing.T, r *redis.Client) {
				userJSON, err := json.Marshal(user)
				require.NoError(t, err)

				err = r.Set(ctx, login, userJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrUserNotFound),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			user, err := repo.Get(ctx, tc.login)
			tc.checkRes(t, user, err)
		})
	}

}

func TestRedisRepository_Create(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	password := []byte("1231231")
	login := "user1"
	lastID := models.UserID(1)
	ctx := context.Background()

	user := models.User{
		ID:       lastID,
		Login:    login,
		Password: password,
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	repo := &RedisRepository{
		client: client,
	}

	expErr := func(expErr error) func(t *testing.T, res *models.User, err error) {
		return func(t *testing.T, res *models.User, err error) {
			assert.ErrorIs(t, expErr, err)
			assert.Nil(t, res)
		}
	}
	tests := []struct {
		name     string
		login    string
		password []byte
		before   func(t *testing.T, r *redis.Client)
		checkRes func(t *testing.T, res *models.User, err error)
	}{
		{
			name:     "create user: user already exists",
			login:    login,
			password: password,
			before: func(t *testing.T, r *redis.Client) {
				userJSON, err := json.Marshal(user)
				require.NoError(t, err)

				err = r.Set(ctx, login, userJSON, 0).Err()
				require.NoError(t, err)
			},
			checkRes: expErr(ErrUserAlreadyExists),
		}, {
			name:     "create user successfully",
			login:    login,
			password: password,
			before:   func(t *testing.T, r *redis.Client) {},
			checkRes: func(t *testing.T, res *models.User, err error) {
				assert.Equal(t, res, &user)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := client.FlushAll(ctx).Err()
			require.NoError(t, err)

			tc.before(t, client)

			user, err := repo.Create(ctx, tc.login, tc.password)
			tc.checkRes(t, user, err)
		})
	}
}
