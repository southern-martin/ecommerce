package database

import (
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CartBackup is the PostgreSQL model for cart backup persistence.
// This table is created for future use in cart persistence/recovery.
type CartBackup struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"uniqueIndex;not null"`
	CartData  string    `gorm:"type:jsonb;not null"` // JSON-encoded cart
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName returns the table name for CartBackup.
func (CartBackup) TableName() string {
	return "cart_backups"
}

// NewPostgresDB creates a new PostgreSQL connection and auto-migrates the schema.
func NewPostgresDB(dsn string, log zerolog.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto-migrate the cart backup table
	if err := db.AutoMigrate(&CartBackup{}); err != nil {
		log.Warn().Err(err).Msg("failed to auto-migrate cart_backups table")
	}

	log.Info().Msg("PostgreSQL connected and migrated")
	return db, nil
}
