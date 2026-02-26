package database

import (
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/promotion/internal/adapter/postgres"
	"github.com/southern-martin/ecommerce/services/promotion/internal/infrastructure/config"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB creates a new GORM database connection and runs auto-migration.
func NewPostgresDB(cfg config.PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(gormPostgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	log.Info().Msg("connected to PostgreSQL")

	// Auto-migrate tables
	err = db.AutoMigrate(
		&postgres.CouponModel{},
		&postgres.CouponUsageModel{},
		&postgres.FlashSaleModel{},
		&postgres.FlashSaleItemModel{},
		&postgres.BundleModel{},
	)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("database migration completed")

	return db, nil
}
