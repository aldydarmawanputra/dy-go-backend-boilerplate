package database

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-backend-boilerplate/internal/config"
)

var dbNamePattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

func EnsureDatabase(cfg *config.Config) error {
	if !dbNamePattern.MatchString(cfg.DBName) {
		return fmt.Errorf("invalid database name %q", cfg.DBName)
	}

	adminDSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/postgres?sslmode=%s",
		url.QueryEscape(cfg.DBUser),
		url.QueryEscape(cfg.DBPassword),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	admin, err := gorm.Open(postgres.Open(adminDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	sqlDB, err := admin.DB()
	if err != nil {
		return err
	}
	defer func() { _ = sqlDB.Close() }()
	sqlDB.SetConnMaxLifetime(time.Minute)

	var count int64
	if err := admin.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", cfg.DBName).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return admin.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName)).Error
}
