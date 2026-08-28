package database

import (
	"gorm.io/gorm"

	usermod "go-backend-boilerplate/internal/modules/user"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermod.User{},
		&usermod.UserDetail{},
	)
}
