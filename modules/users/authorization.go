package users

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IlnurShafikov/wallet/models"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var ErrAuthorizationFailed = errors.New("authorization failed")

type AuthorizationHandler struct {
	service Service
	log     *zerolog.Logger
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserID models.UserID `json:"user_id"`
}

func RegisterAuthorizationHandler(
	router fiber.Router,
	service Service,
	logger *zerolog.Logger,
) {
	auth := &AuthorizationHandler{
		service: service,
		log:     logger,
	}

	router.Post("/login", auth.authorization)
}

func (h *AuthorizationHandler) authorization(fCtx *fiber.Ctx) error {
	req := loginRequest{}
	if err := json.Unmarshal(fCtx.Body(), &req); err != nil {
		h.log.Err(err).Msg("unmarshal failed")
		return err
	}

	userID, err := h.service.Authorization(context.Background(), req.Login, req.Password)
	if err != nil {
		h.log.Err(err).Msg("authorization failed")
		return err
	}

	h.log.Debug().
		Int("userID", int(userID)).
		Msg("authorization successful")

	err = fCtx.Status(fiber.StatusOK).JSON(loginResponse{
		UserID: userID,
	})

	return err
}
