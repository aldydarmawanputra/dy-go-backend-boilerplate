package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"go-backend-boilerplate/internal/config"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusPaid    Status = "paid"
	StatusFailed  Status = "failed"
	StatusExpired Status = "expired"
)

type ChargeRequest struct {
	OrderID       string
	Amount        int64
	Currency      string
	CustomerEmail string
	Description   string
}

type Charge struct {
	Provider   string
	ExternalID string
	PaymentURL string
	Status     Status
}

type WebhookEvent struct {
	OrderID string
	Status  Status
}

// Gateway is the provider-agnostic payment boundary. Implement it for Midtrans,
// Xendit, Stripe, etc. and wire it in New().
type Gateway interface {
	Charge(ctx context.Context, req ChargeRequest) (*Charge, error)
	ParseWebhook(payload []byte) (*WebhookEvent, error)
}

func New(cfg *config.Config) (Gateway, error) {
	switch cfg.PaymentProvider {
	case "stub", "":
		return stub{}, nil
	default:
		return nil, fmt.Errorf("payment provider %q not implemented", cfg.PaymentProvider)
	}
}

// stub is a no-network placeholder: it fabricates a pending charge and parses a
// simple webhook body. Replace with a real provider implementation.
type stub struct{}

func (stub) Charge(_ context.Context, req ChargeRequest) (*Charge, error) {
	return &Charge{
		Provider:   "stub",
		ExternalID: "stub_" + req.OrderID,
		PaymentURL: "https://payments.example.com/pay/" + req.OrderID,
		Status:     StatusPending,
	}, nil
}

func (stub) ParseWebhook(payload []byte) (*WebhookEvent, error) {
	var e struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, err
	}
	return &WebhookEvent{OrderID: e.OrderID, Status: Status(e.Status)}, nil
}
