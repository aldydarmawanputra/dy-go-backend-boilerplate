package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-backend-boilerplate/internal/config"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	logMode := logger.Info
	if cfg.IsProduction() {
		logMode = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
