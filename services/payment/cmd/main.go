package main

import (
	"context"
	"encoding/json"
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

	"github.com/southern-martin/ecommerce/pkg/metrics"
	"github.com/southern-martin/ecommerce/pkg/tracing"

	grpcAdapter "github.com/southern-martin/ecommerce/services/payment/internal/adapter/grpc"
	httpAdapter "github.com/southern-martin/ecommerce/services/payment/internal/adapter/http"
	"github.com/southern-martin/ecommerce/services/payment/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/payment/internal/domain"
	"github.com/southern-martin/ecommerce/services/payment/internal/infrastructure/config"
	"github.com/southern-martin/ecommerce/services/payment/internal/infrastructure/database"
	natsInfra "github.com/southern-martin/ecommerce/services/payment/internal/infrastructure/nats"
	stripeInfra "github.com/southern-martin/ecommerce/services/payment/internal/infrastructure/stripe"
	"github.com/southern-martin/ecommerce/services/payment/internal/usecase"
)

func main() {
	// Configure logger.
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level.
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Init tracer
	tracerShutdown, tracerErr := tracing.InitTracer("payment-service", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), os.Getenv("ENVIRONMENT"))
	if tracerErr != nil {
		log.Warn().Err(tracerErr).Msg("failed to init tracer")
	} else {
		defer tracerShutdown(context.Background())
	}

	log.Info().Msg("Starting Payment Service")

	// Connect to PostgreSQL.
	db, err := database.NewPostgresDB(cfg.PostgresDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Auto-migrate database models.
	if err := db.AutoMigrate(
		&postgres.PaymentModel{},
		&postgres.SellerWalletModel{},
		&postgres.WalletTransactionModel{},
		&postgres.PayoutModel{},
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}
	log.Info().Msg("Database migration completed")

	// Connect to NATS.
	publisher, err := natsInfra.NewPublisher(cfg.NatsURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	defer publisher.Close()

	// Initialize Stripe client (mock for dev).
	stripeClient := stripeInfra.NewMockStripeClient()

	// Initialize repositories.
	paymentRepo := postgres.NewPaymentRepo(db)
	walletRepo := postgres.NewWalletRepo(db)
	payoutRepo := postgres.NewPayoutRepo(db)

	// Initialize use cases.
	createPaymentUC := usecase.NewCreatePaymentUseCase(paymentRepo, stripeClient, publisher)
	confirmPaymentUC := usecase.NewConfirmPaymentUseCase(paymentRepo, walletRepo, publisher, cfg.PlatformCommissionRate)
	walletUC := usecase.NewWalletUseCase(walletRepo)
	payoutUC := usecase.NewPayoutUseCase(payoutRepo, walletRepo, stripeClient)
	refundUC := usecase.NewRefundUseCase(paymentRepo, walletRepo, stripeClient, publisher)

	// Subscribe to order.created events.
	subscribeOrderCreated(publisher, paymentRepo)

	// Initialize HTTP handler and router.
	handler := httpAdapter.NewHandler(
		paymentRepo,
		createPaymentUC,
		confirmPaymentUC,
		walletUC,
		payoutUC,
		refundUC,
	)
	router := httpAdapter.NewRouter(handler)

	// Start gRPC server.
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			tracing.GRPCUnaryInterceptor(),
			metrics.GRPCUnaryInterceptor("payment-service"),
		),
	)
	paymentGRPC := grpcAdapter.NewPaymentGRPCServer(paymentRepo, refundUC)
	grpcAdapter.RegisterPaymentService(grpcServer, paymentGRPC)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatal().Err(err).Str("port", cfg.GRPCPort).Msg("Failed to listen for gRPC")
	}

	go func() {
		log.Info().Str("port", cfg.GRPCPort).Msg("gRPC server started")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Start HTTP server.
	go func() {
		addr := fmt.Sprintf(":%s", cfg.HTTPPort)
		log.Info().Str("port", cfg.HTTPPort).Msg("HTTP server started")
		if err := router.Run(addr); err != nil {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down Payment Service")

	grpcServer.GracefulStop()

	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	log.Info().Msg("Payment Service stopped")
}

// subscribeOrderCreated subscribes to order.created events and creates pending payment records.
func subscribeOrderCreated(publisher *natsInfra.Publisher, paymentRepo domain.PaymentRepository) {
	_, err := publisher.Subscribe(domain.EventOrderCreated, func(data []byte) {
		var event domain.OrderCreatedEvent
		if err := json.Unmarshal(data, &event); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal order.created event")
			return
		}

		log.Info().Str("order_id", event.OrderID).Msg("Received order.created event")

		payment := &domain.Payment{
			ID:          uuid.New().String(),
			OrderID:     event.OrderID,
			BuyerID:     event.BuyerID,
			AmountCents: event.AmountCents,
			Currency:    event.Currency,
			Status:      domain.PaymentStatusPending,
			Method:      domain.PaymentMethodCard,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := paymentRepo.Create(ctx, payment); err != nil {
			log.Error().Err(err).Str("order_id", event.OrderID).Msg("Failed to create payment from order event")
			return
		}

		log.Info().Str("payment_id", payment.ID).Str("order_id", event.OrderID).Msg("Created pending payment from order event")
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to subscribe to order.created")
	}
}
