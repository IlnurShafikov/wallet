package main

import (
	"github.com/IlnurShafikov/wallet/configs"
	"github.com/IlnurShafikov/wallet/services/auth"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/IlnurShafikov/wallet/services/wallet"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

const appName = "APP"

func errorHandler(fCtx *fiber.Ctx, err error) error {
	return fCtx.Status(http.StatusBadRequest).
		Send([]byte(`{"message":"` + err.Error() + `"}`))
}

func main() {
	cfg, err := configs.Parse()
	if err != nil {
		panic(err)
	}

	userWallet := wallet.NewInMemoryRepository()
	usersRepository := users.NewInMemoryRepository()

	fApp := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		ErrorHandler: errorHandler,
		AppName:      appName,
	})

	_ = wallet.NewHandler(fApp, userWallet)
	_ = auth.NewAuthorization(fApp, usersRepository)

	err = fApp.Listen(cfg.GetServerPort())
	if err != nil {
		panic(err)
	}

}
