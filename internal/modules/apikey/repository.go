package apikey

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, k *APIKey) error
	FindByHash(ctx context.Context, keyHash string) (*APIKey, error)
	TouchLastUsed(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, k *APIKey) error {
	return r.db.WithContext(ctx).Create(k).Error
}

func (r *repository) FindByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	var k APIKey
	err := r.db.WithContext(ctx).First(&k, "key_hash = ?", keyHash).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

func (r *repository) TouchLastUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&APIKey{}).Where("id = ?", id).Update("last_used_at", now).Error
}
