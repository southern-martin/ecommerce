package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"ecommerce/catalog-service/internal/adapter/httpapi"
	"ecommerce/catalog-service/internal/infra/memory"
	"ecommerce/catalog-service/internal/infra/postgres"
	"ecommerce/catalog-service/internal/port"
	"ecommerce/catalog-service/internal/usecase"
)

func main() {
	categoryRepo, attributeRepo, productRepo, cleanup, err := buildRepositories()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	categoryService := usecase.NewCategoryService(categoryRepo)
	attributeService := usecase.NewAttributeService(attributeRepo, categoryRepo)
	productService := usecase.NewProductService(productRepo, categoryRepo, attributeRepo)
	variantService := usecase.NewVariantGenerationService(productRepo, attributeRepo)

	api := httpapi.NewServer(
		categoryService,
		attributeService,
		productService,
		variantService,
		productRepo,
	)

	srv := &http.Server{
		Addr:              ":" + envOrDefault("PORT", "8080"),
		Handler:           api.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("catalog-api listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func buildRepositories() (
	port.CategoryRepository,
	port.AttributeRepository,
	port.ProductRepository,
	func(),
	error,
) {
	mode := strings.ToLower(envOrDefault("REPOSITORY_MODE", "memory"))
	switch mode {
	case "memory":
		store := memory.NewStore()
		return memory.NewCategoryRepo(store), memory.NewAttributeRepo(store), memory.NewProductRepo(store), func() {}, nil
	case "postgres":
		dsn := os.Getenv("DATABASE_URL")
		if strings.TrimSpace(dsn) == "" {
			return nil, nil, nil, nil, fmt.Errorf("DATABASE_URL is required when REPOSITORY_MODE=postgres")
		}
		db, err := postgres.Open(dsn)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		if envBool("AUTO_MIGRATE", true) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			migrationsDir := envOrDefault("MIGRATIONS_DIR", "db/migrations")
			if err := postgres.RunMigrations(ctx, db, migrationsDir); err != nil {
				_ = db.Close()
				return nil, nil, nil, nil, err
			}
		}
		return postgres.NewCategoryRepo(db), postgres.NewAttributeRepo(db), postgres.NewProductRepo(db), func() {
			_ = db.Close()
		}, nil
	default:
		return nil, nil, nil, nil, fmt.Errorf("unsupported REPOSITORY_MODE: %s", mode)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch v {
	case "", "default":
		return fallback
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
