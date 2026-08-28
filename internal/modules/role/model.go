package role

type Role struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;not null" json:"name"`
}

func (Role) TableName() string { return "roles" }

type UserRole struct {
	UserID string `gorm:"type:uuid;primaryKey" json:"user_id"`
	RoleID int    `gorm:"primaryKey" json:"role_id"`
}

func (UserRole) TableName() string { return "user_roles" }
