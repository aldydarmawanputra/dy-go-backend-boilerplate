package user

import "time"

type DetailInput struct {
	Phone     string `json:"phone" validate:"omitempty,max=30"`
	Address   string `json:"address" validate:"omitempty,max=255"`
	City      string `json:"city" validate:"omitempty,max=100"`
	Country   string `json:"country" validate:"omitempty,max=100"`
	Bio       string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url,max=255"`
}

type CreateRequest struct {
	Email    string       `json:"email" validate:"required,email"`
	Name     string       `json:"name" validate:"required,min=2,max=100"`
	Password string       `json:"password" validate:"required,min=6,max=72"`
	Detail   *DetailInput `json:"detail" validate:"omitempty"`
}

type ReplaceRequest struct {
	Name   string       `json:"name" validate:"required,min=2,max=100"`
	Detail *DetailInput `json:"detail" validate:"omitempty"`
}

type PatchDetailInput struct {
	Phone     *string `json:"phone" validate:"omitempty,max=30"`
	Address   *string `json:"address" validate:"omitempty,max=255"`
	City      *string `json:"city" validate:"omitempty,max=100"`
	Country   *string `json:"country" validate:"omitempty,max=100"`
	Bio       *string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url,max=255"`
}

type PatchRequest struct {
	Name   *string           `json:"name" validate:"omitempty,min=2,max=100"`
	Detail *PatchDetailInput `json:"detail" validate:"omitempty"`
}

type DetailResponse struct {
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Country   string `json:"country"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type Response struct {
	ID        string          `json:"id"`
	Email     string          `json:"email"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Detail    *DetailResponse `json:"detail,omitempty"`
}

func ToResponse(u *User) Response {
	r := Response{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.Detail != nil {
		r.Detail = &DetailResponse{
			Phone:     u.Detail.Phone,
			Address:   u.Detail.Address,
			City:      u.Detail.City,
			Country:   u.Detail.Country,
			Bio:       u.Detail.Bio,
			AvatarURL: u.Detail.AvatarURL,
		}
	}
	return r
}

func ToResponseList(users []User) []Response {
	out := make([]Response, 0, len(users))
	for i := range users {
		out = append(out, ToResponse(&users[i]))
	}
	return out
}
