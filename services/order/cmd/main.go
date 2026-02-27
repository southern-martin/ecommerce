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

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	grpcAdapter "github.com/southern-martin/ecommerce/services/order/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/order/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/order/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/order/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/order/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/order/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/order/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	setupLogger(cfg.LogLevel)

	// Init tracer
	tracerShutdown, err := tracing.InitTracer("order-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to init tracer")
	} else {
		defer tracerShutdown(context.Background())
	}

	log.Info().Msg("starting order service")

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// Initialize NATS publisher
	publisher, err := natsInfra.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Initialize repositories
	orderRepo := postgres.NewOrderRepo(db)
	sellerOrderRepo := postgres.NewSellerOrderRepo(db)

	// Initialize use cases
	createOrderUC := usecase.NewCreateOrderUseCase(orderRepo, sellerOrderRepo, publisher)
	getOrderUC := usecase.NewGetOrderUseCase(orderRepo, sellerOrderRepo)
	updateStatusUC := usecase.NewUpdateOrderStatusUseCase(orderRepo, sellerOrderRepo, publisher)
	cancelOrderUC := usecase.NewCancelOrderUseCase(orderRepo, sellerOrderRepo, publisher)

	// Initialize HTTP handler and router
	handler := httpAdapter.NewHandler(createOrderUC, getOrderUC, updateStatusUC, cancelOrderUC)
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
			metrics.GRPCUnaryInterceptor("order-service"),
		),
	)
	grpcSrv := grpcAdapter.NewServer(getOrderUC, updateStatusUC)
	grpcAdapter.RegisterOrderServiceServer(grpcServer, grpcSrv)

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

	log.Info().Msg("shutting down order service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	log.Info().Msg("order service stopped")
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
