package main

import (
	"log"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"github.com/southern-martin/ecommerce/pkg/events"
	"github.com/southern-martin/ecommerce/pkg/logger"
	"github.com/southern-martin/ecommerce/pkg/server"

	usergrpc "github.com/southern-martin/ecommerce/services/user/internal/adapter/grpc"
	userhttp "github.com/southern-martin/ecommerce/services/user/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/user/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/user/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/user/internal/infrastructure/database"
	usernats "github.com/southern-martin/ecommerce/services/user/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/user/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	l := logger.New(cfg.LogLevel, "user-service")

	// Connect to Postgres
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	l.Info().Msg("connected to postgres")

	// Connect to NATS
	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	l.Info().Msg("connected to NATS")

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JetStream context: %v", err)
	}

	// Create publisher
	publisher := usernats.NewPublisher(js)

	// Create repositories
	profileRepo := postgres.NewProfileRepository(db)
	addressRepo := postgres.NewAddressRepository(db)
	sellerRepo := postgres.NewSellerRepository(db)
	followRepo := postgres.NewFollowRepository(db)

	// Create use cases
	profileUC := usecase.NewProfileUseCase(profileRepo, l)
	addressUC := usecase.NewAddressUseCase(addressRepo, l)
	sellerUC := usecase.NewSellerUseCase(sellerRepo, publisher, l)
	followUC := usecase.NewFollowUseCase(followRepo, l)

	// Start NATS subscriber for user.registered events
	sub := events.NewSubscriber(js)
	if err := usernats.StartSubscriber(sub, profileUC, l); err != nil {
		log.Fatalf("failed to start NATS subscriber: %v", err)
	}
	l.Info().Msg("NATS subscriber started for user.registered events")

	// Setup HTTP router
	handler := userhttp.NewHandler(profileUC, addressUC, sellerUC, followUC)
	router := userhttp.NewRouter(handler)

	// Setup gRPC server
	grpcServer := grpc.NewServer()
	userServiceServer := usergrpc.NewUserServiceServer(profileUC, sellerUC)
	usergrpc.RegisterServer(grpcServer, userServiceServer)

	// Start servers with graceful shutdown
	srv := server.New(l)
	srv.Run(cfg.HTTPPort, cfg.GRPCPort, router, grpcServer)
}
