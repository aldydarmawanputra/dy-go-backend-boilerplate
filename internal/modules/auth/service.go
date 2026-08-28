package auth

import (
	"context"
	"errors"
	"html"

	"go-backend-boilerplate/internal/mailer"
	"go-backend-boilerplate/internal/modules/role"
	"go-backend-boilerplate/internal/modules/user"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/hash"
	"go-backend-boilerplate/internal/shared/jwtutil"
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
}

type service struct {
	users   user.Service
	repo    user.Repository
	roles   role.Repository
	jwt     *jwtutil.Manager
	refresh *RefreshStore
	mailer  mailer.Mailer
	worker  *worker.Pool
}

func NewService(users user.Service, repo user.Repository, roles role.Repository, jwt *jwtutil.Manager, refresh *RefreshStore, mail mailer.Mailer, pool *worker.Pool) Service {
	return &service{users: users, repo: repo, roles: roles, jwt: jwt, refresh: refresh, mailer: mail, worker: pool}
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

	email, name := u.Email, u.Name
	s.worker.Submit("welcome-email", func(ctx context.Context) error {
		body := "<p>Welcome, " + html.EscapeString(name) + "!</p>"
		return s.mailer.Send(ctx, email, "Welcome", body)
	})

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

	// Rotation: revoke the used token before issuing a new pair.
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
