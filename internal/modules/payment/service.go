package payment

import (
	"context"

	"github.com/google/uuid"

	pay "go-backend-boilerplate/internal/payment"
	"go-backend-boilerplate/internal/shared/apperror"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest, currency string) (*Payment, error)
	HandleWebhook(ctx context.Context, payload []byte) error
}

type service struct {
	repo    Repository
	gateway pay.Gateway
}

func NewService(repo Repository, gateway pay.Gateway) Service {
	return &service{repo: repo, gateway: gateway}
}

func (s *service) Create(ctx context.Context, req CreateRequest, currency string) (*Payment, error) {
	orderID := uuid.NewString()

	charge, err := s.gateway.Charge(ctx, pay.ChargeRequest{
		OrderID:       orderID,
		Amount:        req.Amount,
		Currency:      currency,
		CustomerEmail: req.CustomerEmail,
		Description:   req.Description,
	})
	if err != nil {
		return nil, apperror.Internal("failed to create charge")
	}

	p := &Payment{
		OrderID:    orderID,
		Provider:   charge.Provider,
		ExternalID: charge.ExternalID,
		Amount:     req.Amount,
		Currency:   currency,
		Status:     string(charge.Status),
		PaymentURL: charge.PaymentURL,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) HandleWebhook(ctx context.Context, payload []byte) error {
	event, err := s.gateway.ParseWebhook(payload)
	if err != nil {
		return apperror.BadRequest("invalid webhook payload")
	}
	if event.OrderID == "" {
		return apperror.BadRequest("missing order id")
	}
	return s.repo.UpdateStatus(ctx, event.OrderID, string(event.Status))
}
