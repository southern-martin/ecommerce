package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	googlegrpc "google.golang.org/grpc"

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	aigrpc "github.com/southern-martin/ecommerce/services/ai/internal/adapter/grpc"
	handler "github.com/southern-martin/ecommerce/services/ai/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/ai/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/aiclient"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/database"
	natspub "github.com/southern-martin/ecommerce/services/ai/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/ai/internal/usecase"
)

func main() {
	// Load config
	cfg := config.Load()

	// Setup logger
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("starting AI service")

	tracerShutdown, err := tracing.InitTracer("ai-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
	if err != nil {
		log.Warn().Err(err).Msg("failed to init tracer")
	} else {
		defer tracerShutdown(context.Background())
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&postgres.EmbeddingModel{},
		&postgres.RecommendationModel{},
		&postgres.AIConversationModel{},
		&postgres.GeneratedContentModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	log.Info().Msg("database migrations completed")

	// Connect to NATS
	publisher, err := natspub.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Create AI client
	aiClient := aiclient.NewMockAIClient(cfg.AIPythonServiceURL)

	// Initialize repositories
	embeddingRepo := postgres.NewEmbeddingRepo(db)
	recommendationRepo := postgres.NewRecommendationRepo(db)
	conversationRepo := postgres.NewAIConversationRepo(db)
	contentRepo := postgres.NewGeneratedContentRepo(db)

	// Initialize use cases
	embeddingUC := usecase.NewEmbeddingUseCase(embeddingRepo, aiClient, publisher)
	recommendationUC := usecase.NewRecommendationUseCase(recommendationRepo, publisher)
	chatbotUC := usecase.NewChatbotUseCase(conversationRepo, aiClient)
	contentUC := usecase.NewContentUseCase(contentRepo, aiClient, publisher)

	// Setup HTTP server
	h := handler.NewHandler(embeddingUC, recommendationUC, chatbotUC, contentUC)
	router := handler.NewRouter(h)

	// Start HTTP server
	go func() {
		log.Info().Str("port", cfg.HTTPPort).Msg("starting HTTP server")
		if err := router.Run(":" + cfg.HTTPPort); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Setup gRPC server
	grpcServer := googlegrpc.NewServer(
		googlegrpc.ChainUnaryInterceptor(
			tracing.GRPCUnaryInterceptor(),
			metrics.GRPCUnaryInterceptor("ai-service"),
		),
	)
	aiGRPCServer := aigrpc.NewServer(embeddingUC, recommendationUC, contentUC)
	aigrpc.RegisterAIServiceServer(grpcServer, aiGRPCServer)

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to listen for gRPC")
		}
		log.Info().Str("port", cfg.GRPCPort).Msg("starting gRPC server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down AI service")
	grpcServer.GracefulStop()
	_ = context.Background()
	log.Info().Msg("AI service stopped")
}
