package role

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	AssignByName(ctx context.Context, userID, roleName string) error
	NamesForUser(ctx context.Context, userID string) ([]string, error)
	List(ctx context.Context) ([]Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AssignByName(ctx context.Context, userID, roleName string) error {
	const q = `
		INSERT INTO user_roles (user_id, role_id)
		SELECT $1, id FROM roles WHERE name = $2
		ON CONFLICT DO NOTHING
	`
	return r.db.WithContext(ctx).Exec(q, userID, roleName).Error
}

func (r *repository) NamesForUser(ctx context.Context, userID string) ([]string, error) {
	const q = `
		SELECT r.name
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`
	var names []string
	if err := r.db.WithContext(ctx).Raw(q, userID).Scan(&names).Error; err != nil {
		return nil, err
	}
	return names, nil
}

func (r *repository) List(ctx context.Context) ([]Role, error) {
	var roles []Role
	if err := r.db.WithContext(ctx).Order("id").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
