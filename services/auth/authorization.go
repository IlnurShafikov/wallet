package auth

import (
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
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
	log            *zerolog.Logger
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID models.UserID `json:"user_id"`
}

func NewAuthorization(
	router fiber.Router,
	userRepository usersGetter,
	hashed hasherPasswordVerify,
	logger *zerolog.Logger,
) *AuthorizationHandler {
	auth := &AuthorizationHandler{
		userRepository: userRepository,
		hashedVerify:   hashed,
		log:            logger,
	}

	router.Post("/login", auth.Authorization)

	return auth
}

func (h *AuthorizationHandler) Authorization(fCtx *fiber.Ctx) error {
	req := loginRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).Msg("error read body")
		return err
	}

	user, err := h.userRepository.Get(req.Login)
	if err != nil {
		h.log.Err(err).Msg("get user failed")

		if errors.Is(err, users.ErrUserNotFound) {
			err = ErrAuthorizationFailed
		}

		return err
	}

	err = h.hashedVerify.Verify(req.Password, user.Password)
	if err != nil {
		h.log.Err(err).Msg("authorization failed")
		return ErrAuthorizationFailed
	}

	h.log.Info().
		Int("userID", int(user.ID)).
		Msg("authorization successful")

	err = fCtx.Status(fiber.StatusOK).JSON(loginResponse{
		UserID: user.ID,
	})

	return err
}
