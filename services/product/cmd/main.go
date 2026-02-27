package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	googlegrpc "google.golang.org/grpc"

	"github.com/nats-io/nats.go"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	"github.com/southern-martin/ecommerce/services/product/internal/adapter/grpc"
	producthttp "github.com/southern-martin/ecommerce/services/product/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/product/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/product/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/product/internal/infrastructure/database"
	natspub "github.com/southern-martin/ecommerce/services/product/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/product/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Configure logger
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Init tracer
	tracerShutdown, err := tracing.InitTracer("product-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to init tracer")
	} else {
		defer tracerShutdown(context.Background())
	}

	log.Info().Msg("Starting product service...")

	// Connect to PostgreSQL
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	// Auto-migrate database tables
	if err := db.AutoMigrate(
		&postgres.ProductModel{},
		&postgres.CategoryModel{},
		&postgres.AttributeDefinitionModel{},
		&postgres.CategoryAttributeModel{},
		&postgres.ProductAttributeValueModel{},
		&postgres.ProductOptionModel{},
		&postgres.ProductOptionValueModel{},
		&postgres.VariantModel{},
		&postgres.VariantOptionValueModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to auto-migrate database")
	}
	log.Info().Msg("Database migration completed")

	// Connect to NATS (publisher — plain NATS)
	publisher, err := natspub.NewPublisher(cfg.NatsURL)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to NATS, events will not be published")
	}

	// Connect to NATS JetStream for subscribing to order events
	var jsSub *events.Subscriber
	nc, err := nats.Connect(cfg.NatsURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to NATS for subscriber")
	} else {
		defer nc.Close()
		js, jsErr := nc.JetStream()
		if jsErr != nil {
			log.Warn().Err(jsErr).Msg("Failed to create JetStream context")
		} else {
			jsSub = events.NewSubscriber(js)
			log.Info().Msg("NATS JetStream subscriber context ready")
		}
	}

	// Initialize repositories
	productRepo := postgres.NewProductRepo(db)
	categoryRepo := postgres.NewCategoryRepo(db)
	attributeRepo := postgres.NewAttributeRepo(db)
	optionRepo := postgres.NewOptionRepo(db)
	variantRepo := postgres.NewVariantRepo(db)

	// Initialize use cases
	productUC := usecase.NewProductUseCase(productRepo, categoryRepo, attributeRepo, optionRepo, variantRepo, publisher)
	categoryUC := usecase.NewCategoryUseCase(categoryRepo)
	attributeUC := usecase.NewAttributeUseCase(attributeRepo, categoryRepo)
	variantUC := usecase.NewVariantUseCase(productRepo, optionRepo, variantRepo, publisher)

	// Start NATS subscriber for order.created events (stock decrement)
	if jsSub != nil {
		if err := natspub.StartSubscriber(jsSub, variantUC, log.Logger); err != nil {
			log.Warn().Err(err).Msg("Failed to start order.created subscriber, stock will not be decremented automatically")
		} else {
			log.Info().Msg("NATS subscriber started for order.created events")
		}
	}

	// Initialize HTTP handler and router
	handler := producthttp.NewHandler(productUC, categoryUC, attributeUC, variantUC)
	router := producthttp.NewRouter(handler)

	// Start HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.HTTPPort).Msg("HTTP server listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer(productUC, variantUC)
	go func() {
		if err := grpcServer.Start(cfg.GRPCPort,
			googlegrpc.ChainUnaryInterceptor(
				tracing.GRPCUnaryInterceptor(),
				metrics.GRPCUnaryInterceptor("product-service"),
			),
		); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down product service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	grpcServer.Stop()

	if publisher != nil {
		publisher.Close()
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	log.Info().Msg("Product service stopped")
}
