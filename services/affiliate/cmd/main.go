package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	grpcAdapter "github.com/southern-martin/ecommerce/services/affiliate/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/affiliate/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/affiliate/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/usecase"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	log.Info().Msg("starting affiliate service")

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&postgres.AffiliateProgramModel{},
		&postgres.AffiliateLinkModel{},
		&postgres.ReferralModel{},
		&postgres.AffiliatePayoutModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to auto-migrate")
	}

	// Seed default program
	seedDefaultProgram(db)

	// Initialize NATS publisher
	publisher, err := natsInfra.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Initialize repositories
	programRepo := postgres.NewProgramRepo(db)
	linkRepo := postgres.NewLinkRepo(db)
	referralRepo := postgres.NewReferralRepo(db)
	payoutRepo := postgres.NewPayoutRepo(db)

	// Initialize use cases
	programUC := usecase.NewProgramUseCase(programRepo)
	linkUC := usecase.NewLinkUseCase(linkRepo, publisher)
	referralUC := usecase.NewReferralUseCase(referralRepo, linkRepo, programRepo, publisher)
	payoutUC := usecase.NewPayoutUseCase(payoutRepo, programRepo, linkRepo, publisher)

	// Initialize HTTP handler and router
	handler := httpAdapter.NewHandler(programUC, linkUC, referralUC, payoutUC)
	router := httpAdapter.NewRouter(handler)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler: router,
	}

	go func() {
		log.Info().Str("port", cfg.HTTPPort).Msg("HTTP server starting")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer()
	grpcSrv := grpcAdapter.NewServer(linkUC, referralUC, programUC)
	grpcAdapter.RegisterAffiliateServiceServer(grpcServer, grpcSrv)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen for gRPC")
		}
		log.Info().Str("port", cfg.GRPCPort).Msg("gRPC server starting")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down affiliate service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	log.Info().Msg("affiliate service stopped")
}

func setupLogger(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
}

func seedDefaultProgram(db *gorm.DB) {
	var count int64
	db.Model(&postgres.AffiliateProgramModel{}).Count(&count)
	if count > 0 {
		log.Info().Msg("affiliate program already exists, skipping seed")
		return
	}

	program := postgres.AffiliateProgramModel{
		ID:                 uuid.New().String(),
		CommissionRate:     0.05,
		MinPayoutCents:     5000,
		CookieDays:         30,
		ReferrerBonusCents: 500,
		ReferredBonusCents: 500,
		IsActive:           true,
	}

	if err := db.Create(&program).Error; err != nil {
		log.Warn().Err(err).Msg("failed to seed default affiliate program")
		return
	}

	log.Info().Msg("default affiliate program seeded")
}
