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

	grpcAdapter "github.com/southern-martin/ecommerce/services/media/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/media/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/media/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/media/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/media/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/media/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/media/internal/infrastructure/storage"
	"github.com/southern-martin/ecommerce/services/media/internal/usecase"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	log.Info().Msg("starting media service")

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&postgres.MediaModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to auto-migrate")
	}

	// Initialize NATS publisher
	publisher, err := natsInfra.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Initialize storage client (mock for now)
	storageClient := storage.NewMockStorageClient(cfg.S3.Endpoint, cfg.S3.Bucket)

	// Initialize repositories
	mediaRepo := postgres.NewMediaRepo(db)

	// Initialize use cases
	mediaUC := usecase.NewMediaUseCase(mediaRepo, storageClient, publisher)

	// Initialize HTTP handler and router
	handler := httpAdapter.NewHandler(mediaUC)
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
	grpcSrv := grpcAdapter.NewServer(mediaUC)
	grpcAdapter.RegisterMediaServiceServer(grpcServer, grpcSrv)

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

	log.Info().Msg("shutting down media service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	log.Info().Msg("media service stopped")
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
