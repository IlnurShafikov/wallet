package main

import (
	"github.com/IlnurShafikov/wallet/configs"
	"github.com/IlnurShafikov/wallet/services/auth"
	"github.com/IlnurShafikov/wallet/services/transaction"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/IlnurShafikov/wallet/services/wallet"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

const (
	appName = "APP"
)

func errorHandler(fCtx *fiber.Ctx, err error) error {
	return fCtx.Status(http.StatusBadRequest).
		Send([]byte(`{"message":"` + err.Error() + `"}`))
}

func main() {
	cfg, err := configs.Parse()
	if err != nil {
		panic(err)
	}

	err = cfg.Validate()
	if err != nil {
		panic(err)
	}

	loggerLevelStr := cfg.LogLevel
	loggerLevel, err := zerolog.ParseLevel(loggerLevelStr)
	if err != nil {
		loggerLevel = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Logger().
		Level(loggerLevel)

	walletTransaction := transaction.NewInMemoryRepository()
	usersRepository := users.NewInMemoryRepository()
	hasherPassword := auth.NewBcryptHashing(cfg.Secret)
	userWallet := wallet.NewInMemoryRepository()
	walletTR := wallet.NewWallet(userWallet, walletTransaction)

	fApp := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		ErrorHandler: errorHandler,
		AppName:      appName,
	})

	_ = wallet.NewHandler(fApp, walletTR, &logger)
	_ = auth.NewAuthorization(fApp, usersRepository, hasherPassword, &logger)
	_ = auth.NewRegistrationHandler(fApp, usersRepository, hasherPassword, &logger)
	_ = transaction.NewHandler(fApp, walletTransaction, &logger)

	err = fApp.Listen(cfg.GetServerPort())
	if err != nil {
		panic(err)
	}

}
