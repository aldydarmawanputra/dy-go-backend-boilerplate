package auth

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log/slog"
	"time"

	"go-backend-boilerplate/internal/config"
	"go-backend-boilerplate/internal/mailer"
	"go-backend-boilerplate/internal/modules/role"
	"go-backend-boilerplate/internal/modules/user"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/hash"
	"go-backend-boilerplate/internal/shared/jwtutil"
	"go-backend-boilerplate/internal/shared/redact"
	"go-backend-boilerplate/internal/shared/sanitize"
	"go-backend-boilerplate/internal/worker"
)

type Tokens struct {
	Access  string
	Refresh string
}

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*user.User, error)
	Login(ctx context.Context, req LoginRequest) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
	VerifyEmail(ctx context.Context, token string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword redact.Secret) error
}

type service struct {
	cfg     *config.Config
	users   user.Service
	repo    user.Repository
	roles   role.Repository
	jwt     *jwtutil.Manager
	refresh *RefreshStore
	tokens  *TokenStore
	mailer  mailer.Mailer
	worker  *worker.Pool
}

func NewService(cfg *config.Config, users user.Service, repo user.Repository, roles role.Repository, jwt *jwtutil.Manager, refresh *RefreshStore, tokens *TokenStore, mail mailer.Mailer, pool *worker.Pool) Service {
	return &service{cfg: cfg, users: users, repo: repo, roles: roles, jwt: jwt, refresh: refresh, tokens: tokens, mailer: mail, worker: pool}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*user.User, error) {
	u, err := s.users.Create(ctx, user.CreateRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	s.sendVerificationEmail(ctx, u.ID, u.Email, u.Name)
	return u, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*Tokens, error) {
	u, err := s.repo.FindByEmail(ctx, sanitize.Email(req.Email))
	if err != nil {
		return nil, err
	}
	if u == nil {
		hash.DummyCompare(req.Password.Reveal())
		return nil, apperror.Unauthorized("invalid email or password")
	}
	if err := hash.Compare(u.PasswordHash, req.Password.Reveal()); err != nil {
		return nil, apperror.Unauthorized("invalid email or password")
	}
	return s.issueTokens(ctx, u.ID)
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	userID, err := s.refresh.UserID(ctx, refreshToken)
	if errors.Is(err, ErrRefreshUnavailable) {
		return nil, apperror.New(503, "REFRESH_UNAVAILABLE", "refresh is not available")
	}
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, apperror.Unauthorized("invalid or expired refresh token")
	}
	_ = s.refresh.Revoke(ctx, refreshToken)
	return s.issueTokens(ctx, userID)
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	err := s.refresh.Revoke(ctx, refreshToken)
	if errors.Is(err, ErrRefreshUnavailable) {
		return nil
	}
	return err
}

func (s *service) VerifyEmail(ctx context.Context, token string) error {
	userID, err := s.tokens.Consume(ctx, nsVerify, token)
	if errors.Is(err, ErrTokenStoreUnavailable) {
		return apperror.New(503, "TOKEN_STORE_UNAVAILABLE", "verification is not available")
	}
	if err != nil {
		return err
	}
	if userID == "" {
		return apperror.BadRequest("invalid or expired token")
	}
	return s.repo.MarkEmailVerified(ctx, userID)
}

func (s *service) ForgotPassword(ctx context.Context, email string) error {
	u, err := s.repo.FindByEmail(ctx, sanitize.Email(email))
	if err != nil {
		return err
	}
	// Always succeed to avoid leaking which emails are registered.
	if u == nil {
		return nil
	}
	token, err := s.tokens.Issue(ctx, nsReset, u.ID, time.Duration(s.cfg.ResetTokenExpireHours)*time.Hour)
	if err != nil {
		slog.Warn("issue reset token", "err", err)
		return nil
	}

	link := s.cfg.AppBaseURL + "/reset-password?token=" + token
	email2, name := u.Email, u.Name
	s.worker.Submit("reset-email", func(ctx context.Context) error {
		body := fmt.Sprintf("<p>Hi %s,</p><p>Reset your password: <a href=\"%s\">%s</a></p>",
			html.EscapeString(name), link, link)
		return s.mailer.Send(ctx, email2, "Reset your password", body)
	})
	return nil
}

func (s *service) ResetPassword(ctx context.Context, token string, newPassword redact.Secret) error {
	userID, err := s.tokens.Consume(ctx, nsReset, token)
	if errors.Is(err, ErrTokenStoreUnavailable) {
		return apperror.New(503, "TOKEN_STORE_UNAVAILABLE", "reset is not available")
	}
	if err != nil {
		return err
	}
	if userID == "" {
		return apperror.BadRequest("invalid or expired token")
	}
	hashed, err := hash.Password(newPassword.Reveal())
	if err != nil {
		return apperror.Internal("failed to hash password")
	}
	return s.repo.UpdatePassword(ctx, userID, hashed)
}

func (s *service) sendVerificationEmail(ctx context.Context, userID, email, name string) {
	token, err := s.tokens.Issue(ctx, nsVerify, userID, time.Duration(s.cfg.VerifyTokenExpireHours)*time.Hour)
	if err != nil {
		if !errors.Is(err, ErrTokenStoreUnavailable) {
			slog.Warn("issue verify token", "err", err)
		}
		return
	}
	link := s.cfg.AppBaseURL + "/verify-email?token=" + token
	s.worker.Submit("verify-email", func(ctx context.Context) error {
		body := fmt.Sprintf("<p>Hi %s,</p><p>Verify your email: <a href=\"%s\">%s</a></p>",
			html.EscapeString(name), link, link)
		return s.mailer.Send(ctx, email, "Verify your email", body)
	})
}

func (s *service) issueTokens(ctx context.Context, userID string) (*Tokens, error) {
	roles, err := s.roles.NamesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	access, err := s.jwt.Generate(userID, roles)
	if err != nil {
		return nil, apperror.Internal("failed to issue token")
	}
	refresh, err := s.refresh.Issue(ctx, userID)
	if err != nil && !errors.Is(err, ErrRefreshUnavailable) {
		return nil, apperror.Internal("failed to issue refresh token")
	}
	return &Tokens{Access: access, Refresh: refresh}, nil
}
