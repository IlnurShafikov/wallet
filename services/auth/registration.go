package auth

import (
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/gofiber/fiber/v2"
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

type hasherPassword interface {
	HashPassword(password string) ([]byte, error)
}

type RegistrationHandler struct {
	userCreate userCreater
	hashed     hasherPassword
}

func NewRegistrationHandler(
	router fiber.Router,
	userCreate userCreater,
	hashed hasherPassword,
) *RegistrationHandler {
	handler := &RegistrationHandler{
		userCreate: userCreate,
		hashed:     hashed,
	}

	router.Post("/registration", handler.Registration)

	return handler
}

func (c *RegistrationHandler) Registration(fCtx *fiber.Ctx) error {
	req := CreateUserRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		return err
	}

	if req.Password != req.RePassword {
		return ErrWrongRePassword
	}

	hashPassword, err := c.hashed.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user, err := c.userCreate.Create(req.Login, hashPassword)
	if err != nil {
		return err
	}

	err = fCtx.Status(fiber.StatusCreated).JSON(CreateUserResponse{
		Login:  req.Login,
		UserID: user.ID,
	})

	return err
}
