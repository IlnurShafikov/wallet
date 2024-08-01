package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var (
	ErrWrongRePassword = errors.New("passwords not equal")
)

type CreateUserRequest struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
}

type CreateUserResponse struct {
	Login  string        `json:"login"`
	UserID models.UserID `json:"user_id"`
}

type passwordHasher interface {
	HashPassword(password string) ([]byte, error)
}

type RegistrationHandler struct {
	usersCreater users.Creater
	hasher       passwordHasher
	log          *zerolog.Logger
}

func RunRegistrationHandler(
	router fiber.Router,
	usersCreater users.Creater,
	hashed passwordHasher,
	logger *zerolog.Logger,
) {
	handler := &RegistrationHandler{
		usersCreater: usersCreater,
		hasher:       hashed,
		log:          logger,
	}

	router.Post("/registration", handler.registration)
}

func (c *RegistrationHandler) registration(fCtx *fiber.Ctx) error {
	req := CreateUserRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		c.log.Err(err).Msg("unmarshal failed")
		return err
	}

	if req.Password != req.RePassword {
		c.log.Warn().Msg("password not equal")
		return ErrWrongRePassword
	}

	hashPassword, err := c.hasher.HashPassword(req.Password)
	if err != nil {
		c.log.Err(err).Msg("hashing password failed")
		return err
	}
	ctx := context.Background()

	user, err := c.usersCreater.Create(ctx, req.Login, hashPassword)
	if err != nil {
		c.log.Err(err).Msg("registration failed")
		return err
	}

	c.log.Debug().
		Int("userID", int(user.ID)).
		Msg("registration successful")

	err = fCtx.Status(fiber.StatusCreated).JSON(CreateUserResponse{
		Login:  req.Login,
		UserID: user.ID,
	})

	return err
}
