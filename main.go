package main

import (
	"context"
	"fmt"
	"github.com/IlnurShafikov/wallet/configs"
	"github.com/IlnurShafikov/wallet/services/auth"
	"github.com/IlnurShafikov/wallet/services/auth/security"
	"github.com/IlnurShafikov/wallet/services/transaction"
	"github.com/IlnurShafikov/wallet/services/users"
	"github.com/IlnurShafikov/wallet/services/wallet"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"
	"time"
)

type components struct {
	userRepository        users.Repository
	walletRepository      wallet.Repository
	transactionRepository transaction.Repository
}

const (
	appName = "APP"
)

func errorHandler(fCtx *fiber.Ctx, err error) error {
	return fCtx.Status(http.StatusBadRequest).
		JSON(struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		})
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := configs.Parse()
	if err != nil {
		return err
	}

	err = cfg.Validate()
	if err != nil {
		return err
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

	comp, err := makeComponents(cfg)
	if err != nil {
		return fmt.Errorf("failed create components: %w", err)
	}

	hasherPassword := security.NewBcryptHashing(cfg.Secret)
	walletTR := wallet.NewWallet(comp.walletRepository, comp.transactionRepository, &logger)

	fApp := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		ErrorHandler: errorHandler,
		AppName:      appName,
	})

	wallet.RunWalletHandler(fApp, walletTR, &logger)
	auth.RunAuthorizationHandler(fApp, comp.userRepository, hasherPassword, &logger)
	auth.RunRegistrationHandler(fApp, comp.userRepository, hasherPassword, &logger)
	transaction.RunTransactionHandler(fApp, comp.transactionRepository, &logger)

	err = fApp.Listen(cfg.GetServerPort())
	if err != nil {
		return err
	}

	return nil
}

func makeComponents(cfg *configs.Config) (*components, error) {
	switch cfg.StorageType {
	case "in_memory":
		return inMemoryComponent()
	case "redis":
		return redisComponent(cfg)
	default:
		return nil, fmt.Errorf("unknow storage type: %s", cfg.StorageType)
	}
}

func redisComponent(cfg *configs.Config) (*components, error) {
	clientRedis := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Address,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := clientRedis.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis.Ping:%w", err)
	}

	resp := &components{
		userRepository:        users.NewRedisRepository(clientRedis),
		walletRepository:      wallet.NewRedisRepository(clientRedis),
		transactionRepository: transaction.NewRedisRepository(clientRedis),
	}

	return resp, nil
}

func inMemoryComponent() (*components, error) {
	resp := &components{
		userRepository:        users.NewInMemoryRepository(),
		walletRepository:      wallet.NewInMemoryRepository(),
		transactionRepository: transaction.NewInMemoryRepository(),
	}

	return resp, nil
}
