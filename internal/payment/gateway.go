package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"go-backend-boilerplate/internal/config"
)

var ErrInvalidSignature = errors.New("invalid webhook signature")

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
	// ParseWebhook verifies the signature over the raw payload, then decodes it.
	ParseWebhook(payload []byte, signature string) (*WebhookEvent, error)
}

func New(cfg *config.Config) (Gateway, error) {
	switch cfg.PaymentProvider {
	case "stub", "":
		return stub{secret: cfg.PaymentWebhookSecret}, nil
	default:
		return nil, fmt.Errorf("payment provider %q not implemented", cfg.PaymentProvider)
	}
}

// stub is a no-network placeholder: it fabricates a pending charge and parses a
// simple webhook body. Replace with a real provider implementation.
type stub struct {
	secret string
}

func (stub) Charge(_ context.Context, req ChargeRequest) (*Charge, error) {
	return &Charge{
		Provider:   "stub",
		ExternalID: "stub_" + req.OrderID,
		PaymentURL: "https://payments.example.com/pay/" + req.OrderID,
		Status:     StatusPending,
	}, nil
}

func (s stub) ParseWebhook(payload []byte, signature string) (*WebhookEvent, error) {
	if s.secret != "" && !validHMAC(s.secret, payload, signature) {
		return nil, ErrInvalidSignature
	}
	var e struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
	}
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, err
	}
	return &WebhookEvent{OrderID: e.OrderID, Status: Status(e.Status)}, nil
}

func validHMAC(secret string, payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
