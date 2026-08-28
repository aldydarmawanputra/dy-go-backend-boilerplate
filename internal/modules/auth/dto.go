package auth

import "go-backend-boilerplate/internal/shared/redact"

type RegisterRequest struct {
	Email    string        `json:"email" validate:"required,email"`
	Name     string        `json:"name" validate:"required,min=2,max=100"`
	Password redact.Secret `json:"password" validate:"required,min=6,max=72"`
}

type LoginRequest struct {
	Email    string        `json:"email" validate:"required,email"`
	Password redact.Secret `json:"password" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string        `json:"token" validate:"required"`
	Password redact.Secret `json:"password" validate:"required,min=6,max=72"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
}

func toTokenResponse(t *Tokens) TokenResponse {
	return TokenResponse{
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
		TokenType:    "Bearer",
	}
}
