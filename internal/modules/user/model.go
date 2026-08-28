package user

import "time"

type User struct {
	ID           string      `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string      `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string      `gorm:"not null" json:"-"`
	Name         string      `gorm:"not null" json:"name"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Detail       *UserDetail `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"detail,omitempty"`
}

func (User) TableName() string { return "users" }

type UserDetail struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserDetail) TableName() string { return "user_details" }
