package transaction

import (
	"context"
	"github.com/IlnurShafikov/wallet/models"
)

type Repository interface {
	GetRound(context.Context, models.RoundID) (*models.Round, error)
	CreateBet(context.Context, models.RoundID, models.Round) error
	SetWin(context.Context, models.RoundID, models.Transaction) error
	UpdateRound(context.Context, models.RoundID, models.Round) error
}
