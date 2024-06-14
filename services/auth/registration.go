package auth

import (
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
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

type userCreater interface {
	Create(login string, password []byte) (*models.User, error)
}

type passwordHasher interface {
	HashPassword(password string) ([]byte, error)
}

type RegistrationHandler struct {
	userCreate userCreater
	hasher     passwordHasher
	log        *zerolog.Logger
}

func NewRegistrationHandler(
	router fiber.Router,
	userCreate userCreater,
	hashed passwordHasher,
	logger *zerolog.Logger,
) *RegistrationHandler {
	handler := &RegistrationHandler{
		userCreate: userCreate,
		hasher:     hashed,
		log:        logger,
	}

	router.Post("/registration", handler.Registration)

	return handler
}

func (c *RegistrationHandler) Registration(fCtx *fiber.Ctx) error {
	req := CreateUserRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		c.log.Err(err).Msg("error read body")
		return err
	}

	if req.Password != req.RePassword {
		c.log.Err(ErrWrongRePassword).Msg("password not equal")
		return ErrWrongRePassword
	}

	hashPassword, err := c.hasher.HashPassword(req.Password)
	if err != nil {
		c.log.Err(err).Msg("hashing password failed")
		return err
	}

	user, err := c.userCreate.Create(req.Login, hashPassword)
	if err != nil {
		c.log.Err(err).Msg("registration failed")
		return err
	}

	c.log.Info().
		Int("userID", int(user.ID)).
		Msg("registration successful")

	err = fCtx.Status(fiber.StatusCreated).JSON(CreateUserResponse{
		Login:  req.Login,
		UserID: user.ID,
	})

	return err
}
