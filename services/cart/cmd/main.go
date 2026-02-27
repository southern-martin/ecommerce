package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	"github.com/southern-martin/ecommerce/services/cart/internal/adapter/dualstore"
	cartgrpc "github.com/southern-martin/ecommerce/services/cart/internal/adapter/grpc"
	carthttp "github.com/southern-martin/ecommerce/services/cart/internal/adapter/http"
	cartpg "github.com/southern-martin/ecommerce/services/cart/internal/adapter/postgres"
	cartredis "github.com/southern-martin/ecommerce/services/cart/internal/adapter/redis"
	"github.com/southern-martin/ecommerce/services/cart/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/cart/internal/infrastructure/database"
	cartnats "github.com/southern-martin/ecommerce/services/cart/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/cart/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("service", "cart").Logger().Level(level)

	// Init tracer
	tracerShutdown, tracerErr := tracing.InitTracer("cart-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
	if tracerErr != nil {
		logger.Warn().Err(tracerErr).Msg("failed to init tracer")
	} else {
		defer tracerShutdown(context.Background())
	}

	// Connect to PostgreSQL (durable cart persistence)
	pgDB, err := database.NewPostgresDB(cfg.PostgresDSN(), logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to PostgreSQL")
	}

	// Connect to Redis (fast cache layer)
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to Redis")
	}
	logger.Info().Str("addr", cfg.RedisURL).Msg("Redis connected")
	defer rdb.Close()

	// Connect to NATS
	natsConn, err := cartnats.Connect(cfg.NATSURL, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer natsConn.Close()

	// Initialize layers: Redis cache + Postgres durable store
	redisRepo := cartredis.NewRedisCartRepository(rdb)
	pgRepo := cartpg.NewPostgresCartRepository(pgDB)
	cartRepo := dualstore.NewDualStoreCartRepository(redisRepo, pgRepo, logger)
	eventPublisher := cartnats.NewEventPublisher(natsConn, logger)
	cartUC := usecase.NewCartUseCase(cartRepo, eventPublisher, logger)

	// HTTP server
	handler := carthttp.NewCartHandler(cartUC, logger)
	router := carthttp.NewRouter(handler)
	httpServer := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			tracing.GRPCUnaryInterceptor(),
			metrics.GRPCUnaryInterceptor("cart-service"),
		),
	)
	cartSvc := cartgrpc.NewCartServiceServer(cartUC, logger)
	cartgrpc.RegisterCartServiceServer(grpcServer, cartSvc)

	grpcLis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		logger.Fatal().Err(err).Str("port", cfg.GRPCPort).Msg("failed to listen for gRPC")
	}

	// Start servers
	go func() {
		logger.Info().Str("port", cfg.HTTPPort).Msg("HTTP server starting")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	go func() {
		logger.Info().Str("port", cfg.GRPCPort).Msg("gRPC server starting")
		if err := grpcServer.Serve(grpcLis); err != nil {
			logger.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info().Str("signal", sig.String()).Msg("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("HTTP server shutdown error")
	}

	logger.Info().Msg("cart service stopped")
}
