package payment

type CreateRequest struct {
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Description   string `json:"description" validate:"omitempty,max=255"`
	CustomerEmail string `json:"customer_email" validate:"omitempty,email"`
}
