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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	grpcAdapter "github.com/southern-martin/ecommerce/services/loyalty/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/loyalty/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/loyalty/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/loyalty/internal/usecase"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	log.Info().Msg("starting loyalty service")

	tracerShutdown, err := tracing.InitTracer("loyalty-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to init tracer")
	} else {
		defer tracerShutdown(context.Background())
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&postgres.MembershipModel{},
		&postgres.PointsTransactionModel{},
		&postgres.TierModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to auto-migrate")
	}

	// Seed default tiers
	seedTiers(db)

	// Initialize NATS publisher
	publisher, err := natsInfra.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Initialize repositories
	membershipRepo := postgres.NewMembershipRepo(db)
	transactionRepo := postgres.NewTransactionRepo(db)
	tierRepo := postgres.NewTierRepo(db)

	// Initialize use cases
	membershipUC := usecase.NewMembershipUseCase(membershipRepo, tierRepo, publisher)
	pointsUC := usecase.NewPointsUseCase(membershipRepo, transactionRepo, membershipUC, publisher)
	tierUC := usecase.NewTierUseCase(tierRepo)

	// Initialize HTTP handler and router
	handler := httpAdapter.NewHandler(membershipUC, pointsUC, tierUC)
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
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			tracing.GRPCUnaryInterceptor(),
			metrics.GRPCUnaryInterceptor("loyalty-service"),
		),
	)
	grpcSrv := grpcAdapter.NewServer(membershipUC, pointsUC)
	grpcAdapter.RegisterLoyaltyServiceServer(grpcServer, grpcSrv)

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

	log.Info().Msg("shutting down loyalty service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	log.Info().Msg("loyalty service stopped")
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

func seedTiers(db *gorm.DB) {
	defaultTiers := []postgres.TierModel{
		{Name: "bronze", MinPoints: 0, CashbackRate: 0.01, PointsMultiplier: 1.0, FreeShipping: false, PrioritySupportHours: 48},
		{Name: "silver", MinPoints: 1000, CashbackRate: 0.02, PointsMultiplier: 1.2, FreeShipping: false, PrioritySupportHours: 24},
		{Name: "gold", MinPoints: 5000, CashbackRate: 0.03, PointsMultiplier: 1.5, FreeShipping: true, PrioritySupportHours: 12},
		{Name: "platinum", MinPoints: 15000, CashbackRate: 0.05, PointsMultiplier: 2.0, FreeShipping: true, PrioritySupportHours: 4},
	}

	for _, tier := range defaultTiers {
		result := db.Where("name = ?", tier.Name).FirstOrCreate(&tier)
		if result.Error != nil {
			log.Warn().Err(result.Error).Str("name", tier.Name).Msg("failed to seed tier")
		}
	}

	log.Info().Msg("default tiers seeded")
}
