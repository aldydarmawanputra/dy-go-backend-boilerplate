package database

import (
	"gorm.io/gorm"

	paymentmod "go-backend-boilerplate/internal/modules/payment"
	rolemod "go-backend-boilerplate/internal/modules/role"
	usermod "go-backend-boilerplate/internal/modules/user"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&usermod.User{},
		&usermod.UserDetail{},
		&rolemod.Role{},
		&rolemod.UserRole{},
		&paymentmod.Payment{},
	)
}
