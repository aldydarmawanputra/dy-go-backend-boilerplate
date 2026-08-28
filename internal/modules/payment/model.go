package payment

import "go-backend-boilerplate/internal/shared/model"

type Payment struct {
	model.Base
	OrderID    string `gorm:"uniqueIndex;not null" json:"order_id"`
	Provider   string `gorm:"not null" json:"provider"`
	ExternalID string `json:"external_id"`
	Amount     int64  `gorm:"not null" json:"amount"`
	Currency   string `gorm:"not null" json:"currency"`
	Status     string `gorm:"not null" json:"status"`
	PaymentURL string `json:"payment_url"`
}

func (Payment) TableName() string { return "payments" }
