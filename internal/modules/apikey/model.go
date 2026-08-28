package apikey

import (
	"time"

	"go-backend-boilerplate/internal/shared/model"
)

type APIKey struct {
	model.Base
	Name       string     `gorm:"not null" json:"name"`
	KeyHash    string     `gorm:"uniqueIndex;not null" json:"-"`
	Prefix     string     `gorm:"not null" json:"prefix"`
	Revoked    bool       `gorm:"not null;default:false" json:"revoked"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

func (APIKey) TableName() string { return "api_keys" }
