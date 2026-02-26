package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/southern-martin/ecommerce/services/tax/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/tax/internal/domain"
	"github.com/southern-martin/ecommerce/services/tax/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/tax/internal/infrastructure/database"
	"github.com/southern-martin/ecommerce/services/tax/internal/usecase"

	grpcAdapter "github.com/southern-martin/ecommerce/services/tax/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/tax/internal/adapter/http"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg := config.Load()

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err == nil {
		zerolog.SetGlobalLevel(level)
	}

	log.Info().Msg("Starting Tax Service")

	// Connect to database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(&postgres.TaxZoneModel{}, &postgres.TaxRuleModel{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	log.Info().Msg("Database migrations completed")

	// Create repositories
	zoneRepo := postgres.NewTaxZoneRepository(db)
	ruleRepo := postgres.NewTaxRuleRepository(db)

	// Seed default data
	seedDefaultData(zoneRepo, ruleRepo)

	// Create use cases
	calculateTaxUC := usecase.NewCalculateTaxUseCase(zoneRepo, ruleRepo)
	manageRulesUC := usecase.NewManageRulesUseCase(ruleRepo)
	manageZonesUC := usecase.NewManageZonesUseCase(zoneRepo)

	// Setup HTTP server
	handler := httpAdapter.NewHandler(calculateTaxUC, manageRulesUC, manageZonesUC)
	router := httpAdapter.NewRouter(handler)

	// Setup gRPC server
	grpcServer := grpc.NewServer()
	taxGRPCServer := grpcAdapter.NewServer(calculateTaxUC, manageRulesUC, manageZonesUC)
	grpcServer.RegisterService(&grpcAdapter.ServiceDesc, taxGRPCServer)

	// Start gRPC server
	grpcAddr := fmt.Sprintf(":%s", cfg.GRPCPort)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", grpcAddr).Msg("Failed to listen for gRPC")
	}

	go func() {
		log.Info().Str("addr", grpcAddr).Msg("gRPC server listening")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Start HTTP server
	httpAddr := fmt.Sprintf(":%s", cfg.HTTPPort)
	go func() {
		log.Info().Str("addr", httpAddr).Msg("HTTP server listening")
		if err := router.Run(httpAddr); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down Tax Service")
	grpcServer.GracefulStop()
	log.Info().Msg("Tax Service stopped")
}

func seedDefaultData(zoneRepo domain.TaxZoneRepository, ruleRepo domain.TaxRuleRepository) {
	ctx := context.Background()

	// Check if zones already exist
	zones, err := zoneRepo.List(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check existing zones")
		return
	}
	if len(zones) > 0 {
		log.Info().Int("count", len(zones)).Msg("Tax zones already seeded, skipping")
		return
	}

	log.Info().Msg("Seeding default tax zones and rules")

	now := time.Now()

	// Define seed zones
	seedZones := []domain.TaxZone{
		{ID: uuid.New().String(), CountryCode: "AU", StateCode: "", Name: "Australia"},
		{ID: uuid.New().String(), CountryCode: "US", StateCode: "CA", Name: "California"},
		{ID: uuid.New().String(), CountryCode: "US", StateCode: "NY", Name: "New York"},
		{ID: uuid.New().String(), CountryCode: "GB", StateCode: "", Name: "United Kingdom"},
		{ID: uuid.New().String(), CountryCode: "CA", StateCode: "ON", Name: "Ontario, Canada"},
	}

	// Define seed rules (will be linked to zone IDs after creation)
	type seedRule struct {
		zoneIndex int
		rule      domain.TaxRule
	}

	seedRules := []seedRule{
		{
			zoneIndex: 0, // AU
			rule: domain.TaxRule{
				ID: uuid.New().String(), TaxName: "GST", Rate: 0.10,
				Category: "", Inclusive: true, StartsAt: now, IsActive: true,
			},
		},
		{
			zoneIndex: 1, // US-CA
			rule: domain.TaxRule{
				ID: uuid.New().String(), TaxName: "State Sales Tax", Rate: 0.0725,
				Category: "", Inclusive: false, StartsAt: now, IsActive: true,
			},
		},
		{
			zoneIndex: 2, // US-NY
			rule: domain.TaxRule{
				ID: uuid.New().String(), TaxName: "State Sales Tax", Rate: 0.08,
				Category: "", Inclusive: false, StartsAt: now, IsActive: true,
			},
		},
		{
			zoneIndex: 3, // GB
			rule: domain.TaxRule{
				ID: uuid.New().String(), TaxName: "VAT", Rate: 0.20,
				Category: "", Inclusive: true, StartsAt: now, IsActive: true,
			},
		},
		{
			zoneIndex: 4, // CA-ON
			rule: domain.TaxRule{
				ID: uuid.New().String(), TaxName: "HST", Rate: 0.13,
				Category: "", Inclusive: false, StartsAt: now, IsActive: true,
			},
		},
	}

	// Create zones
	for i := range seedZones {
		if err := zoneRepo.Create(ctx, &seedZones[i]); err != nil {
			log.Error().Err(err).Str("zone", seedZones[i].Name).Msg("Failed to seed zone")
			return
		}
	}

	// Create rules with zone IDs
	for _, sr := range seedRules {
		sr.rule.ZoneID = seedZones[sr.zoneIndex].ID
		if err := ruleRepo.Create(ctx, &sr.rule); err != nil {
			log.Error().Err(err).Str("rule", sr.rule.TaxName).Msg("Failed to seed rule")
			return
		}
	}

	log.Info().
		Int("zones", len(seedZones)).
		Int("rules", len(seedRules)).
		Msg("Default tax data seeded successfully")
}
