package user

import "go-backend-boilerplate/internal/shared/model"

type User struct {
	model.Base
	Email         string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string      `gorm:"not null" json:"-"`
	Name          string      `gorm:"not null" json:"name"`
	EmailVerified bool        `gorm:"not null;default:false" json:"email_verified"`
	Detail        *UserDetail `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"detail,omitempty"`
	Roles         []string    `gorm:"-" json:"roles,omitempty"`
}

func (User) TableName() string { return "users" }

type UserDetail struct {
	model.Base
	UserID    string `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

func (UserDetail) TableName() string { return "user_details" }
