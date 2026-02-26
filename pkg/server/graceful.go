package server

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
	"google.golang.org/grpc"
)

// Server manages both an HTTP and a gRPC server with graceful shutdown.
type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	logger     zerolog.Logger
}

// New creates a new Server with the given logger.
func New(logger zerolog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

// Run starts both the HTTP and gRPC servers and handles OS signals for graceful shutdown.
func (s *Server) Run(httpPort, grpcPort string, httpHandler http.Handler, grpcServer *grpc.Server) {
	s.grpcServer = grpcServer
	s.httpServer = &http.Server{
		Addr:         ":" + httpPort,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			s.logger.Fatal().Err(err).Str("port", grpcPort).Msg("failed to listen for gRPC")
		}
		s.logger.Info().Str("port", grpcPort).Msg("gRPC server starting")
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// Start HTTP server
	go func() {
		s.logger.Info().Str("port", httpPort).Msg("HTTP server starting")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	s.waitForShutdown()
}

// waitForShutdown blocks until an OS interrupt or termination signal is received,
// then gracefully shuts down both servers.
func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	s.logger.Info().Str("signal", sig.String()).Msg("shutting down servers")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully stop gRPC
	if s.grpcServer != nil {
		s.logger.Info().Msg("stopping gRPC server")
		s.grpcServer.GracefulStop()
	}

	// Gracefully stop HTTP
	if s.httpServer != nil {
		s.logger.Info().Msg("stopping HTTP server")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error().Err(err).Msg("HTTP server forced to shutdown")
		}
	}

	fmt.Println("servers stopped")
}
