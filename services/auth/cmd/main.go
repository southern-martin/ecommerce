package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"github.com/southern-martin/ecommerce/pkg/auth"
	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/pkg/logger"
	"github.com/southern-martin/ecommerce/pkg/server"

	authgrpc "github.com/southern-martin/ecommerce/services/auth/internal/adapter/grpc"
	authhttp "github.com/southern-martin/ecommerce/services/auth/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/auth/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/auth/internal/domain"
	"github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/database"
	authnats "github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/nats"
	authredis "github.com/southern-martin/ecommerce/services/auth/internal/infrastructure/redis"
	"github.com/southern-martin/ecommerce/services/auth/internal/usecase"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init logger
	l := logger.New(cfg.LogLevel, "auth-service")

	// Connect to Postgres
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// AutoMigrate the AuthUser model
	if err := db.AutoMigrate(&domain.AuthUser{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
	l.Info().Msg("database migration completed")

	// Connect to Redis
	blacklist, err := authredis.NewTokenBlacklist(cfg.RedisURL)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer blacklist.Close()
	l.Info().Msg("connected to redis")

	// Connect to NATS JetStream
	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create jetstream context: %v", err)
	}
	l.Info().Msg("connected to nats jetstream")

	natsPublisher := events.NewNATSPublisher(js)
	eventPublisher := authnats.NewEventPublisher(natsPublisher)

	// Parse JWT expiries
	accessExpiry := auth.ParseExpiry(cfg.JWTAccessExpiry)
	refreshExpiry := auth.ParseExpiry(cfg.JWTRefreshExpiry)

	// Create repository
	repo := postgres.NewUserRepository(db)

	// Create use cases
	registerUC := usecase.NewRegisterUseCase(repo, eventPublisher, cfg.JWTSecret, accessExpiry, refreshExpiry, l)
	loginUC := usecase.NewLoginUseCase(repo, cfg.JWTSecret, accessExpiry, refreshExpiry, l)
	refreshTokenUC := usecase.NewRefreshTokenUseCase(repo, cfg.JWTSecret, accessExpiry, refreshExpiry, l)
	logoutUC := usecase.NewLogoutUseCase(repo, blacklist, cfg.JWTSecret, l)
	forgotPasswordUC := usecase.NewForgotPasswordUseCase(repo, eventPublisher, l)
	resetPasswordUC := usecase.NewResetPasswordUseCase(repo, l)
	oauthLoginUC := usecase.NewOAuthLoginUseCase(repo, eventPublisher, cfg.JWTSecret, accessExpiry, refreshExpiry, l)
	updateRoleUC := usecase.NewUpdateRoleUseCase(repo, l)

	// Setup HTTP handler and router
	handler := authhttp.NewHandler(
		registerUC, loginUC, refreshTokenUC, logoutUC,
		forgotPasswordUC, resetPasswordUC, oauthLoginUC,
	)
	router := authhttp.SetupRouter(handler)

	// Setup gRPC server
	grpcSrv := grpc.NewServer()
	authGRPCServer := authgrpc.NewAuthServiceServer(cfg.JWTSecret, updateRoleUC, l)
	authGRPCServer.Register(grpcSrv)

	// Run both servers with graceful shutdown
	srv := server.New(l)
	srv.Run(cfg.HTTPPort, cfg.GRPCPort, router, grpcSrv)
}
