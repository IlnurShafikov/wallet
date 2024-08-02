package users

import (
	"context"
	"github.com/IlnurShafikov/wallet/models"
)

type Repository interface {
	Creater
	Getter
}

type Creater interface {
	Create(ctx context.Context, login string, password []byte) (*models.User, error)
}

type Getter interface {
	Get(ctx context.Context, login string) (*models.User, error)
}
