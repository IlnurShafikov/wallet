package wallet

import (
	"context"
	"errors"
	"github.com/IlnurShafikov/wallet/mocks"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/modules/wallet/request"
	"github.com/IlnurShafikov/wallet/services/transaction"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type mock struct {
	walletRepo *mocks.MockwalletRepository
	mockTrRepo *mocks.MocktransactionRepository
}

func newMock(ctrl *gomock.Controller) *mock {
	mockWalletRepo := mocks.NewMockwalletRepository(ctrl)
	mockTrRepo := mocks.NewMocktransactionRepository(ctrl)

	return &mock{
		walletRepo: mockWalletRepo,
		mockTrRepo: mockTrRepo,
	}
}

func (m *mock) getWallet(ctx context.Context, userID models.UserID) func(
	balance models.Balance, err error,
) *gomock.Call {
	return func(expBalance models.Balance, expErr error) *gomock.Call {
		return m.walletRepo.EXPECT().
			Get(ctx, userID).Times(1).Return(expBalance, expErr)
	}
}

func (m *mock) getRound(ctx context.Context, roundID models.RoundID) func(
	round *models.Round, err error,
) *gomock.Call {
	return func(expRound *models.Round, expErr error) *gomock.Call {
		return m.mockTrRepo.EXPECT().
			GetRound(ctx, roundID).Times(1).Return(expRound, expErr)
	}
}

func (m *mock) update(ctx context.Context, userID models.UserID, amount models.Amount) func(
	balance models.Balance, err error,
) *gomock.Call {
	return func(expBalance models.Balance, expErr error) *gomock.Call {
		return m.walletRepo.EXPECT().
			Update(ctx, userID, amount).Times(1).Return(expBalance, expErr)
	}
}

func (m *mock) updateRound(ctx context.Context, roundID models.RoundID, round models.Round) func(
	err error,
) *gomock.Call {
	return func(expErr error) *gomock.Call {
		return m.mockTrRepo.EXPECT().
			UpdateRound(ctx, roundID, round).Times(1).Return(expErr)
	}
}

func (m *mock) createBetTr(ctx context.Context, roundID models.RoundID, round models.Round) func(
	err error,
) *gomock.Call {
	return func(expErr error) *gomock.Call {
		return m.mockTrRepo.EXPECT().
			CreateBet(ctx, roundID, round).Times(1).Return(expErr)
	}
}

func (m *mock) setWinTr(ctx context.Context, roundID models.RoundID, winRound models.Transaction) func(
	err error,
) *gomock.Call {
	return func(expErr error) *gomock.Call {
		return m.mockTrRepo.EXPECT().
			SetWin(ctx, roundID, winRound).Times(1).Return(expErr)
	}
}

func TestWallet_Get(t *testing.T) {
	const (
		userID  = 1992
		balance = 141
	)

	ctx := context.Background()

	tests := []struct {
		name   string
		before func(m *mock)
		exp    models.Balance
		err    error
	}{
		{
			name: "wallet not found",
			before: func(m *mock) {
				m.getWallet(ctx, userID)(0, ErrWalletNotFound)
			},
			exp: 0,
			err: ErrWalletNotFound,
		},
		{
			name: "success",
			before: func(m *mock) {
				m.getWallet(ctx, userID)(balance, nil)
			},
			exp: balance,
			err: nil,
		},
	}

	log := zerolog.Nop()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)

			srv := NewWallet(m.walletRepo, m.mockTrRepo, &log)
			balance, err := srv.Get(ctx, userID)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.exp, balance)
		})
	}
}

func TestWallet_Refund(t *testing.T) {
	const (
		userID  = 1992
		amount  = 100
		balance = 1000
	)
	var (
		errRoundNotFound = errors.New("round not found")
	)

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	betID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")
	require.NoError(t, err)

	ctx := context.Background()

	req := request.RefundTransaction{
		RoundID: models.RoundID(roundID),
	}

	log := zerolog.Nop()

	tests := []struct {
		name          string
		before        func(m *mock)
		expectBalance models.Balance
		expectErr     error
	}{
		{
			name: "round not found",
			before: func(m *mock) {
				m.getRound(ctx, models.RoundID(roundID))(nil, errRoundNotFound)
			},
			expectErr: errRoundNotFound,
		}, {
			name: "wallet not found",
			before: func(m *mock) {
				round := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: false,
					Refunded: false,
				}

				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(&round, nil),
					m.update(ctx, userID, amount)(0, ErrWalletNotFound),
				)
			},
			expectBalance: 0,
			expectErr:     ErrWalletNotFound,
		}, {
			name: "refund exists",
			before: func(m *mock) {
				round := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win: &models.Transaction{
						Amount:        amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Finished: true,
					Refunded: false,
				}
				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(&round, nil),
				)
			},
			expectBalance: 0,
			expectErr:     ErrNotRefund,
		}, {
			name: "refund successful",
			before: func(m *mock) {
				round := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: false,
					Refunded: false,
				}

				roundRef := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: false,
					Refunded: true,
				}

				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(&round, nil),
					m.update(ctx, userID, amount)(balance, err),
					m.updateRound(ctx, models.RoundID(roundID), roundRef)(nil),
				)
			},
			expectBalance: balance,
			expectErr:     nil,
		}, {
			name: "update round failed",
			before: func(m *mock) {
				round := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: false,
					Refunded: false,
				}

				roundRef := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: false,
					Refunded: true,
				}

				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(&round, nil),
					m.update(ctx, userID, amount)(balance, err),
					m.updateRound(ctx, models.RoundID(roundID), roundRef)(ErrUpdateRoundFailed),
				)
			},
			expectBalance: 0,
			expectErr:     ErrUpdateRoundFailed,
		},
		{
			name: "refund successful",
			before: func(m *mock) {
				roundRef := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        -amount,
						TransactionID: models.TransactionID(betID),
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: false,
					Refunded: true,
				}

				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(&roundRef, nil),
				)
			},
			expectBalance: 0,
			expectErr:     ErrRefundAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)

			srv := NewWallet(m.walletRepo, m.mockTrRepo, &log)
			balance, err := srv.Refund(ctx, userID, req)

			assert.ErrorIs(t, tt.expectErr, err)
			assert.Equal(t, tt.expectBalance, balance)

		})
	}
}

func TestWallet_CreateBet(t *testing.T) {
	const (
		userID  = 1992
		balance = 100
		amount  = -10
	)

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	betID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")
	require.NoError(t, err)
	ctx := context.Background()

	var (
		errGetFailed    = errors.New("get transaction failed")
		errSetBetFailed = errors.New("set bet failed")
	)
	req := request.UpdateBalance{
		Amount:        amount,
		RoundID:       models.RoundID(roundID),
		TransactionID: models.TransactionID(betID),
		Finished:      false,
	}

	tests := []struct {
		name       string
		before     func(m *mock)
		expBalance models.Balance
		expErr     error
	}{
		{
			name: "refund - update round failed",
			before: func(m *mock) {
				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(nil, transaction.ErrRoundNotFound),
					m.update(ctx, userID, amount)(0, ErrWalletNotEnoughMoney),
				)
			},
			expBalance: 0,
			expErr:     ErrWalletNotEnoughMoney,
		}, {
			name: "create bet successful",
			before: func(m *mock) {
				roundBet := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        req.Amount,
						TransactionID: req.TransactionID,
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: req.Finished,
					Refunded: false,
				}

				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(nil, transaction.ErrRoundNotFound),
					m.update(ctx, userID, amount)(balance, err),
					m.createBetTr(ctx, req.RoundID, roundBet)(nil),
				)
			},
			expBalance: balance,
			expErr:     nil,
		}, {
			name: "create bet transaction failed",
			before: func(m *mock) {
				roundBet := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        req.Amount,
						TransactionID: req.TransactionID,
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: req.Finished,
					Refunded: false,
				}

				gomock.InOrder(
					m.getRound(ctx, models.RoundID(roundID))(nil, transaction.ErrRoundNotFound),
					m.update(ctx, userID, amount)(balance, err),
					m.createBetTr(ctx, req.RoundID, roundBet)(errSetBetFailed),
				)
			},
			expBalance: 0,
			expErr:     errSetBetFailed,
		},
		{
			name: "round already exists",
			before: func(m *mock) {
				roundBet := models.Round{}

				m.getRound(ctx, models.RoundID(roundID))(&roundBet, nil)
			},
			expBalance: 0,
			expErr:     ErrRoundIDAlready,
		},
		{
			name: "get round transaction failed",
			before: func(m *mock) {
				m.getRound(ctx, models.RoundID(roundID))(nil, errGetFailed)
			},
			expBalance: 0,
			expErr:     errGetFailed,
		},
	}

	log := zerolog.Nop()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)

			srv := NewWallet(m.walletRepo, m.mockTrRepo, &log)
			balance, err := srv.createBet(ctx, userID, req)
			assert.ErrorIs(t, err, tt.expErr)
			assert.Equal(t, tt.expBalance, balance)
		})
	}
}

func TestWallet_SetWin(t *testing.T) {
	const (
		userID  = 1992
		balance = 100
		amount  = -10
	)

	var (
		errRoundNotFound           = errors.New("round not found")
		errSetWinTransactionFailed = errors.New("set win transaction failed")
	)

	roundID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	require.NoError(t, err)

	betID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")
	require.NoError(t, err)

	ctx := context.Background()

	req := request.UpdateBalance{
		Amount:        amount,
		RoundID:       models.RoundID(roundID),
		TransactionID: models.TransactionID(betID),
		Finished:      false,
	}

	tests := []struct {
		name       string
		before     func(m *mock)
		expBalance models.Balance
		expErr     error
	}{
		{
			name: "set win successful",
			before: func(m *mock) {
				roundBet := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        req.Amount,
						TransactionID: req.TransactionID,
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: req.Finished,
					Refunded: false,
				}

				winRound := models.Transaction{
					Amount:        req.Amount,
					TransactionID: req.TransactionID,
					Created:       time.Now(),
				}

				gomock.InOrder(
					m.getRound(ctx, req.RoundID)(&roundBet, nil),
					m.update(ctx, userID, req.Amount)(balance, nil),
					m.setWinTr(ctx, req.RoundID, winRound)(nil),
				)
			},
			expBalance: balance,
			expErr:     nil,
		}, {
			name: "round not found",
			before: func(m *mock) {
				m.getRound(ctx, models.RoundID(roundID))(nil, errRoundNotFound)
			},
			expErr: errRoundNotFound,
		},
		{
			name: "win already exists",
			before: func(m *mock) {
				winRound := models.Transaction{
					Amount:        req.Amount,
					TransactionID: req.TransactionID,
					Created:       time.Now(),
				}

				roundBet := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        req.Amount,
						TransactionID: req.TransactionID,
						Created:       time.Now(),
					},
					Win:      &winRound,
					Finished: req.Finished,
					Refunded: false,
				}
				m.getRound(ctx, req.RoundID)(&roundBet, nil)
			},
			expBalance: 0,
			expErr:     ErrWinAlreadyExists,
		}, {
			name: "round finished",
			before: func(m *mock) {
				winRound := models.Transaction{
					Amount:        req.Amount,
					TransactionID: req.TransactionID,
					Created:       time.Now(),
				}

				roundBet := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        req.Amount,
						TransactionID: req.TransactionID,
						Created:       time.Now(),
					},
					Win:      &winRound,
					Finished: true,
					Refunded: false,
				}
				m.getRound(ctx, req.RoundID)(&roundBet, nil)
			},
			expBalance: 0,
			expErr:     ErrRoundFinished,
		}, {
			name: "set win transaction failed",
			before: func(m *mock) {
				roundBet := models.Round{
					UserID: userID,
					Bet: models.Transaction{
						Amount:        req.Amount,
						TransactionID: req.TransactionID,
						Created:       time.Now(),
					},
					Win:      nil,
					Finished: req.Finished,
					Refunded: false,
				}

				winRound := models.Transaction{
					Amount:        req.Amount,
					TransactionID: req.TransactionID,
					Created:       time.Now(),
				}

				gomock.InOrder(
					m.getRound(ctx, req.RoundID)(&roundBet, nil),
					m.update(ctx, userID, req.Amount)(balance, nil),
					m.setWinTr(ctx, req.RoundID, winRound)(errSetWinTransactionFailed),
				)
			},
			expBalance: 0,
			expErr:     errSetWinTransactionFailed,
		},
	}

	log := zerolog.Nop()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := newMock(ctrl)
			tt.before(m)

			srv := NewWallet(m.walletRepo, m.mockTrRepo, &log)
			balance, err := srv.setWin(ctx, userID, req)
			assert.ErrorIs(t, tt.expErr, err)
			assert.Equal(t, tt.expBalance, balance)
		})
	}
}
