package auth

import (
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/gofiber/fiber/v2"
)

var ErrAuthorizationFailed = errors.New("authorization failed")

type usersGetter interface {
	Get(login string) (*models.User, error)
}

type hasherPasswordVerify interface {
	Verify(password string, hashPassword []byte) error
}

type AuthorizationHandler struct {
	userRepository usersGetter
	hashedVerify   hasherPasswordVerify
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID models.UserID `json:"user_id"`
}

func NewAuthorization(
	router fiber.Router,
	userRepository usersGetter,
	hashed hasherPasswordVerify,
) *AuthorizationHandler {
	auth := &AuthorizationHandler{
		userRepository: userRepository,
		hashedVerify:   hashed,
	}

	router.Post("/login", auth.Authorization)

	return auth
}

func (h *AuthorizationHandler) Authorization(fCtx *fiber.Ctx) error {
	req := LoginRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		return err
	}

	user, err := h.userRepository.Get(req.Login)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			err = ErrAuthorizationFailed
		}

		return err
	}

	err = h.hashedVerify.Verify(req.Password, user.Password)
	if err != nil {
		return ErrAuthorizationFailed
	}

	err = fCtx.Status(fiber.StatusOK).JSON(LoginResponse{
		UserID: user.ID,
	})

	return err
}
