package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	secret   []byte
	issuer   string
	audience string
	expire   time.Duration
}

func New(secret, issuer, audience string, expireHours int) *Manager {
	if expireHours <= 0 {
		expireHours = 24
	}
	return &Manager{
		secret:   []byte(secret),
		issuer:   issuer,
		audience: audience,
		expire:   time.Duration(expireHours) * time.Hour,
	}
}

func (m *Manager) Generate(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    m.issuer,
		Subject:   userID,
		Audience:  jwt.ClaimStrings{m.audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(m.expire)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        uuid.NewString(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

func (m *Manager) Parse(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return m.secret, nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(m.audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	if claims.Subject == "" {
		return "", errors.New("invalid token: empty subject")
	}
	return claims.Subject, nil
}
