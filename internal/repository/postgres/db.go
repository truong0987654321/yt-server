package postgres

import (
	"fmt"
	"yt-clone-server/internal/config"
	"yt-clone-server/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	logLevel := logger.Silent

	if cfg.Env == "development" {
		logLevel = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %w", err)
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&domain.Category{},
		&domain.Channel{},
	); err != nil {
		return nil, fmt.Errorf("Failed to run database migrations: %w", err)
	}

	return db, nil
}
