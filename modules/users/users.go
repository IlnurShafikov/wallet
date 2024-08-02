package users

import (
	"context"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/modules/users/repositories"
	"github.com/IlnurShafikov/wallet/services/security"
)

type Service interface {
	Authorization(ctx context.Context, login, password string) (models.UserID, error)
}

type UserService struct {
	repository   Repository
	hashedVerify security.PasswordVerify
}

func NewUserService(
	repository Repository,
	hashedVerify security.PasswordVerify,
) *UserService {
	return &UserService{
		repository:   repository,
		hashedVerify: hashedVerify,
	}
}

func (u *UserService) Authorization(ctx context.Context, login, password string) (models.UserID, error) {
	user, err := u.repository.Get(ctx, login)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			err = ErrAuthorizationFailed
		}
		return 0, err
	}

	err = u.hashedVerify.Verify(password, user.Password)
	if err != nil {
		return 0, ErrAuthorizationFailed
	}

	return user.ID, nil
}
