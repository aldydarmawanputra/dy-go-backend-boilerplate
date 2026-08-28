package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend-boilerplate/internal/shared/audit"
)

type Base struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedBy *string        `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *string        `gorm:"type:uuid" json:"updated_by,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = uuid.NewString()
	}
	if actor := audit.Actor(tx.Statement.Context); actor != "" {
		a := actor
		if b.CreatedBy == nil {
			b.CreatedBy = &a
		}
		b.UpdatedBy = &a
	}
	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) error {
	if actor := audit.Actor(tx.Statement.Context); actor != "" {
		a := actor
		b.UpdatedBy = &a
	}
	return nil
}
