package auth

import (
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/gofiber/fiber/v2"
)

var ErrAuthorizationFailed = errors.New("authorization failed")

type userRepository interface {
	Get(login string) (*models.User, error)
}

type AuthorizationHandler struct {
	userRepository userRepository
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID models.UserID
}

func NewAuthorization(
	router fiber.Router,
	userRepository userRepository,
) *AuthorizationHandler {
	auth := &AuthorizationHandler{userRepository: userRepository}

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

	if string(user.Password) != req.Password {
		return ErrAuthorizationFailed
	}

	err = fCtx.Status(fiber.StatusOK).JSON(LoginResponse{
		UserID: user.ID,
	})

	return err
}
