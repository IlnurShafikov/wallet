package repositories

import (
	"context"
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	const loginUser = "ilnur"
	ctx := context.Background()

	tests := []struct {
		name      string
		login     string
		password  []byte
		before    func(nw *InMemoryRepository)
		expect    *models.User
		expectErr error
	}{
		{
			name:     "create new user",
			login:    loginUser,
			password: []byte("123"),
			before:   func(nw *InMemoryRepository) {},
			expect:   &models.User{1, loginUser, []byte("123")},
		}, {
			name:     "creat an existing user",
			login:    loginUser,
			password: []byte("123"),
			before: func(nw *InMemoryRepository) {
				nw.users[loginUser] = models.User{1, loginUser, []byte("123")}
			},

			expectErr: fmt.Errorf("this user %s exists", loginUser),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nw := NewInMemoryRepository()
			tc.before(nw)
			got, err := nw.Create(ctx, tc.login, tc.password)
			assert.Equal(t, tc.expect, got)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestGet(t *testing.T) {
	password := []byte("123")
	loginUser := "user01"
	loginWrongUser := "user23"
	ctx := context.Background()

	tests := []struct {
		name      string
		login     string
		expect    *models.User
		expectErr error
	}{
		{
			name:   "get real user",
			login:  loginUser,
			expect: &models.User{1, loginUser, password},
		},
		{
			name:      "get wrong user",
			login:     loginWrongUser,
			expectErr: ErrUserNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nw := NewInMemoryRepository()
			_, _ = nw.Create(ctx, loginUser, password)
			got, err := nw.Get(ctx, tc.login)
			assert.Equal(t, tc.expect, got)
			assert.ErrorIs(t, err, tc.expectErr)
		})
	}
}
