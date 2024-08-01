package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/IlnurShafikov/wallet/services/auth/security"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var ErrAuthorizationFailed = errors.New("authorization failed")

type AuthorizationHandler struct {
	userRepository users.Getter
	hashedVerify   security.PasswordVerify
	log            *zerolog.Logger
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID models.UserID `json:"user_id"`
}

func RunAuthorizationHandler(
	router fiber.Router,
	usersGetter users.Getter,
	hashed security.PasswordVerify,
	logger *zerolog.Logger,
) {
	auth := &AuthorizationHandler{
		userRepository: usersGetter,
		hashedVerify:   hashed,
		log:            logger,
	}

	router.Post("/login", auth.authorization)
}

func (h *AuthorizationHandler) authorization(fCtx *fiber.Ctx) error {
	req := loginRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).Msg("unmarshal failed")
		return err
	}

	ctx := context.Background()

	user, err := h.userRepository.Get(ctx, req.Login)
	if err != nil {
		h.log.Err(err).Msg("get user failed")

		if errors.Is(err, users.ErrUserNotFound) {
			err = ErrAuthorizationFailed
		}

		return err
	}

	err = h.hashedVerify.Verify(req.Password, user.Password)
	if err != nil {
		h.log.Warn().
			Int("userID", int(user.ID)).Err(err).
			Msg("verify user password failed")
		return ErrAuthorizationFailed
	}

	h.log.Debug().
		Int("userID", int(user.ID)).
		Msg("authorization successful")

	err = fCtx.Status(fiber.StatusOK).JSON(loginResponse{
		UserID: user.ID,
	})

	return err
}
