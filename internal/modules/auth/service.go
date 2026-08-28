package auth

import (
	"context"

	"go-backend-boilerplate/internal/modules/user"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/hash"
	"go-backend-boilerplate/internal/shared/jwtutil"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*user.User, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
}

type service struct {
	users user.Service
	repo  user.Repository
	jwt   *jwtutil.Manager
}

func NewService(users user.Service, repo user.Repository, jwt *jwtutil.Manager) Service {
	return &service{users: users, repo: repo, jwt: jwt}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*user.User, error) {
	return s.users.Create(ctx, user.CreateRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, error) {
	u, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", apperror.Unauthorized("invalid email or password")
	}
	if err := hash.Compare(u.PasswordHash, req.Password); err != nil {
		return "", apperror.Unauthorized("invalid email or password")
	}
	token, err := s.jwt.Generate(u.ID)
	if err != nil {
		return "", apperror.Internal("failed to issue token")
	}
	return token, nil
}
