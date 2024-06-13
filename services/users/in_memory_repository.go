package users

import (
	"fmt"
	"github.com/IlnurShafikov/wallet/models"
)

type InMemoryRepository struct {
	users  map[string]models.User
	lastID models.UserID
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users:  make(map[string]models.User),
		lastID: 0,
	}
}

func (i *InMemoryRepository) getUser(login string) (models.User, bool) {
	us, ok := i.users[login]
	return us, ok
}

func (i *InMemoryRepository) Create(login string, password []byte) (*models.User, error) {
	if _, exist := i.getUser(login); exist {
		return nil, fmt.Errorf("this user %s exists", login)
	}

	i.lastID++

	user := models.User{
		UserID:   i.lastID,
		Login:    login,
		Password: password,
	}

	i.users[login] = user

	return &user, nil
}

func (i *InMemoryRepository) Get(login string) (*models.User, error) {
	user, exist := i.getUser(login)
	if !exist {
		return nil, fmt.Errorf("this user %s does not exist", login)
	}

	return &user, nil
}
