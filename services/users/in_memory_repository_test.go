package users

import (
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestCunrurency(t *testing.T) {
	t.Parallel()

	const (
		goroutinesCount = 100
		operationsCount = 1000
	)

	repo := NewInMemoryRepository()
	wg := sync.WaitGroup{}
	wg.Add(goroutinesCount * 2)

	for i := 0; i < goroutinesCount; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 1; j <= operationsCount; j++ {
				_, _ = repo.Create(strconv.Itoa(i), []byte(""))
			}
		}(i)
	}

	for i := 0; i < goroutinesCount; i++ {
		go func(i int) {
			defer wg.Done()

			for j := 0; j < operationsCount; j++ {
				_, _ = repo.Get(strconv.Itoa(i))
			}
		}(i)
	}

	wg.Wait()
}

func TestCreate(t *testing.T) {
	const loginUser = "ilnur"
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
			got, err := nw.Create(tc.login, tc.password)
			assert.Equal(t, tc.expect, got)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestGet(t *testing.T) {
	password := []byte("123")
	loginUser := "user01"
	loginWrongUser := "user23"
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
			_, _ = nw.Create(loginUser, password)
			got, err := nw.Get(tc.login)
			assert.Equal(t, tc.expect, got)
			assert.ErrorIs(t, err, tc.expectErr)
		})
	}
}
