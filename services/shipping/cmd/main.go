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

	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	grpcAdapter "github.com/southern-martin/ecommerce/services/shipping/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/shipping/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/shipping/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/shipping/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/shipping/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/shipping/internal/infrastructure/nats"
	"github.com/southern-martin/ecommerce/services/shipping/internal/usecase"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.LogLevel)

	log.Info().Msg("starting shipping service")

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// AutoMigrate
	if err := db.AutoMigrate(
		&postgres.CarrierModel{},
		&postgres.CarrierCredentialModel{},
		&postgres.ShipmentModel{},
		&postgres.ShipmentItemModel{},
		&postgres.TrackingEventModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("failed to auto-migrate")
	}

	// Initialize NATS publisher
	publisher, err := natsInfra.NewPublisher(cfg.NATS.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to NATS")
	}
	defer publisher.Close()

	// Seed default carriers
	seedCarriers(db)

	// Initialize repositories
	carrierRepo := postgres.NewCarrierRepo(db)
	credentialRepo := postgres.NewCredentialRepo(db)
	shipmentRepo := postgres.NewShipmentRepo(db)
	trackingRepo := postgres.NewTrackingEventRepo(db)

	// Initialize use cases
	rateUC := usecase.NewRateUseCase(carrierRepo)
	shipmentUC := usecase.NewShipmentUseCase(shipmentRepo, publisher)
	labelUC := usecase.NewLabelUseCase(shipmentRepo, publisher)
	trackingUC := usecase.NewTrackingUseCase(shipmentRepo, trackingRepo, publisher)
	carrierUC := usecase.NewCarrierUseCase(carrierRepo, credentialRepo)

	// Initialize HTTP handler and router
	handler := httpAdapter.NewHandler(rateUC, shipmentUC, labelUC, trackingUC, carrierUC)
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
	grpcSrv := grpcAdapter.NewServer(rateUC, shipmentUC)
	grpcAdapter.RegisterShippingServiceServer(grpcServer, grpcSrv)

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

	log.Info().Msg("shutting down shipping service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcServer.GracefulStop()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}

	log.Info().Msg("shipping service stopped")
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

func seedCarriers(db *gorm.DB) {
	defaultCarriers := []postgres.CarrierModel{
		{Code: "fedex", Name: "FedEx", IsActive: true, SupportedCountries: pq.StringArray{"US", "CA", "MX", "GB", "AU"}},
		{Code: "ups", Name: "UPS", IsActive: true, SupportedCountries: pq.StringArray{"US", "CA", "MX", "GB", "DE"}},
		{Code: "dhl", Name: "DHL", IsActive: true, SupportedCountries: pq.StringArray{"US", "GB", "DE", "FR", "AU", "JP"}},
		{Code: "usps", Name: "USPS", IsActive: true, SupportedCountries: pq.StringArray{"US"}},
		{Code: "auspost", Name: "Australia Post", IsActive: true, SupportedCountries: pq.StringArray{"AU", "NZ"}},
	}

	for _, carrier := range defaultCarriers {
		result := db.Where("code = ?", carrier.Code).FirstOrCreate(&carrier)
		if result.Error != nil {
			log.Warn().Err(result.Error).Str("code", carrier.Code).Msg("failed to seed carrier")
		}
	}

	log.Info().Msg("default carriers seeded")
}
