package user

import (
	"context"

	"go-backend-boilerplate/internal/modules/role"
	"go-backend-boilerplate/internal/shared/apperror"
	"go-backend-boilerplate/internal/shared/hash"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Search(ctx context.Context, keyword string, limit, offset int) ([]User, int64, error)
	Replace(ctx context.Context, id string, req ReplaceRequest) (*User, error)
	Patch(ctx context.Context, id string, req PatchRequest) (*User, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo  Repository
	roles role.Repository
}

func NewService(repo Repository, roles role.Repository) Service {
	return &service{repo: repo, roles: roles}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*User, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperror.Conflict("email already registered")
	}

	hashed, err := hash.Password(req.Password.Reveal())
	if err != nil {
		return nil, apperror.Internal("failed to hash password")
	}

	u := &User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashed,
	}
	if req.Detail != nil {
		u.Detail = &UserDetail{
			Phone:     req.Detail.Phone,
			Address:   req.Detail.Address,
			City:      req.Detail.City,
			Country:   req.Detail.Country,
			Bio:       req.Detail.Bio,
			AvatarURL: req.Detail.AvatarURL,
		}
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	if err := s.roles.AssignByName(ctx, u.ID, role.Default); err != nil {
		return nil, err
	}
	u.Roles = []string{role.Default}
	return u, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, apperror.NotFound("user not found")
	}
	names, err := s.roles.NamesForUser(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.Roles = names
	return u, nil
}

func (s *service) Search(ctx context.Context, keyword string, limit, offset int) ([]User, int64, error) {
	return s.repo.Search(ctx, keyword, limit, offset)
}

func (s *service) Replace(ctx context.Context, id string, req ReplaceRequest) (*User, error) {
	u, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Name = req.Name
	if req.Detail != nil {
		if u.Detail == nil {
			u.Detail = &UserDetail{UserID: u.ID}
		}
		u.Detail.Phone = req.Detail.Phone
		u.Detail.Address = req.Detail.Address
		u.Detail.City = req.Detail.City
		u.Detail.Country = req.Detail.Country
		u.Detail.Bio = req.Detail.Bio
		u.Detail.AvatarURL = req.Detail.AvatarURL
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Patch(ctx context.Context, id string, req PatchRequest) (*User, error) {
	u, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Detail != nil {
		if u.Detail == nil {
			u.Detail = &UserDetail{UserID: u.ID}
		}
		applyDetailPatch(u.Detail, req.Detail)
	}

	if err := s.repo.Save(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return apperror.NotFound("user not found")
	}
	return s.repo.Delete(ctx, id)
}

func applyDetailPatch(d *UserDetail, p *PatchDetailInput) {
	if p.Phone != nil {
		d.Phone = *p.Phone
	}
	if p.Address != nil {
		d.Address = *p.Address
	}
	if p.City != nil {
		d.City = *p.City
	}
	if p.Country != nil {
		d.Country = *p.Country
	}
	if p.Bio != nil {
		d.Bio = *p.Bio
	}
	if p.AvatarURL != nil {
		d.AvatarURL = *p.AvatarURL
	}
}
