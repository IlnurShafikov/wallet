package users

import (
	"fmt"
)

type InMemoryRepository struct {
	users  map[string]User
	lastID int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users:  make(map[string]User),
		lastID: 0,
	}
}

func (i *InMemoryRepository) getUser(login string) (User, bool) {
	us, ok := i.users[login]
	return us, ok
}

func (i *InMemoryRepository) Create(login string, password []byte) (*User, error) {
	if _, exist := i.getUser(login); exist {
		return nil, fmt.Errorf("this user %s exists", login)
	}

	i.lastID++

	user := User{
		ID:       i.lastID,
		Login:    login,
		Password: password,
	}

	i.users[login] = user

	return &user, nil
}

func (i *InMemoryRepository) Get(login string) (*User, error) {
	user, exist := i.getUser(login)
	if !exist {
		return nil, fmt.Errorf("this user %s does not exist", login)
	}

	return &user, nil
}
