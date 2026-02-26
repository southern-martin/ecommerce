package database

import (
	"github.com/rs/zerolog/log"
	"github.com/southern-martin/ecommerce/services/affiliate/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB creates a new GORM database connection.
func NewPostgresDB(cfg config.PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	log.Info().Str("host", cfg.Host).Str("db", cfg.DBName).Msg("connected to PostgreSQL")
	return db, nil
}
