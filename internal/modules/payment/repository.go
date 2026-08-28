package payment

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, p *Payment) error
	UpdateStatus(ctx context.Context, orderID, status string) error
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *repository) UpdateStatus(ctx context.Context, orderID, status string) error {
	return r.db.WithContext(ctx).
		Model(&Payment{}).
		Where("order_id = ?", orderID).
		Update("status", status).Error
}

func (r *repository) FindByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	var p Payment
	err := r.db.WithContext(ctx).First(&p, "order_id = ?", orderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
