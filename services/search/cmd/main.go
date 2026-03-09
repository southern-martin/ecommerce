// @title          Search Service API
// @version        1.0
// @description    Full-text product search and autocomplete suggestions.
//
// @host           localhost:28085
// @BasePath       /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token (via Kong gateway)

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

	"github.com/nats-io/nats.go"

	"github.com/southern-martin/ecommerce/pkg/cache"
	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"
	esAdapter "github.com/southern-martin/ecommerce/services/search/internal/adapter/elasticsearch"
	grpcAdapter "github.com/southern-martin/ecommerce/services/search/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/search/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/search/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/search/internal/domain"
	"github.com/southern-martin/ecommerce/services/search/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/search/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/search/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/search/internal/usecase"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	log.Info().Msg("starting search service")

	tracerShutdown, err := tracing.InitTracer("search-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
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
		&postgres.SearchIndexModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to auto-migrate")
	}

	// Initialize NATS publisher
	publisher, err := natsInfra.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Initialize Redis cache
	cacheClient := cache.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	defer cacheClient.Close()
	if err := cacheClient.Ping(context.Background()); err != nil {
		log.Warn().Err(err).Msg("Redis not available, caching disabled")
		cacheClient = nil
	}

	// Initialize search repository — prefer Elasticsearch if configured.
	var searchRepo domain.SearchRepository

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL != "" {
		esRepo, esErr := esAdapter.NewESSearchRepo(esURL)
		if esErr != nil {
			log.Warn().Err(esErr).Msg("failed to connect to elasticsearch, falling back to postgres")
			searchRepo = postgres.NewSearchRepo(db)
			log.Info().Msg("using postgres search backend")
		} else {
			searchRepo = esRepo
			log.Info().Str("url", esURL).Msg("using elasticsearch search backend")
		}
	} else {
		searchRepo = postgres.NewSearchRepo(db)
		log.Info().Msg("ELASTICSEARCH_URL not set, using postgres search backend")
	}

	// Initialize use cases
	searchUC := usecase.NewSearchUseCase(searchRepo)
	indexUC := usecase.NewIndexUseCase(searchRepo, publisher)

	// Initialize NATS JetStream subscriber for product events
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Warn().Err(err).Msg("failed to connect to NATS for subscriber (search will work without event subscriptions)")
	} else {
		js, err := nc.JetStream()
		if err != nil {
			log.Warn().Err(err).Msg("failed to create JetStream context for subscriber")
		} else {
			sub := events.NewSubscriber(js)
			if err := natsInfra.StartSubscriber(sub, indexUC, log.Logger); err != nil {
				log.Warn().Err(err).Msg("failed to start NATS subscribers")
			}
		}
		defer nc.Close()
	}

	// Initialize HTTP handler and router
	handler := httpAdapter.NewHandler(searchUC, indexUC, db)
	router := httpAdapter.NewRouter(handler, cacheClient, log.Logger)

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
			metrics.GRPCUnaryInterceptor("search-service"),
		),
	)
	grpcSrv := grpcAdapter.NewServer(searchUC, indexUC)
	grpcAdapter.RegisterSearchServiceServer(grpcServer, grpcSrv)

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

	log.Info().Msg("shutting down search service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	log.Info().Msg("search service stopped")
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
