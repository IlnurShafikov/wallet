package transaction

import (
	"context"
	"github.com/IlnurShafikov/wallet/models"
)

type RedisRepository struct {
}

func (r *RedisRepository) GetRound(ctx context.Context, id models.RoundID) (*models.Round, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisRepository) CreateBet(ctx context.Context, id models.RoundID, round models.Round) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisRepository) SetWin(ctx context.Context, id models.RoundID, transaction models.Transaction) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisRepository) UpdateRound(ctx context.Context, id models.RoundID, round models.Round) error {
	//TODO implement me
	panic("implement me")
}
