package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Search(ctx context.Context, keyword string, limit, offset int) ([]User, int64, error)
	FullTextSearch(ctx context.Context, query string, limit, offset int) ([]User, int64, error)
	Save(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *repository) FindByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Preload("Detail").First(&u, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Preload("Detail").First(&u, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) Search(ctx context.Context, keyword string, limit, offset int) ([]User, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Placeholders ($1..$3) keep these raw queries safe from SQL injection.
	const filter = `($1 = '' OR email ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%')`

	var total int64
	if err := r.db.WithContext(ctx).
		Raw(`SELECT count(*) FROM users WHERE `+filter, keyword).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	const q = `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE ` + filter + `
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var users []User
	if err := r.db.WithContext(ctx).Raw(q, keyword, limit, offset).Scan(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// FullTextSearch uses the generated tsvector column + GIN index (see migration
// 20260829000005) and ranks results by relevance.
func (r *repository) FullTextSearch(ctx context.Context, query string, limit, offset int) ([]User, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	if err := r.db.WithContext(ctx).
		Raw(`SELECT count(*) FROM users WHERE search_vector @@ plainto_tsquery('simple', $1)`, query).
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	const q = `
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users
		WHERE search_vector @@ plainto_tsquery('simple', $1)
		ORDER BY ts_rank(search_vector, plainto_tsquery('simple', $1)) DESC
		LIMIT $2 OFFSET $3
	`
	var users []User
	if err := r.db.WithContext(ctx).Raw(q, query, limit, offset).Scan(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *repository) Save(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(u).Error
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error
}
