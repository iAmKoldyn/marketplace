package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/marketplace/internal/auth"
	"github.com/yourusername/marketplace/internal/store"
)

type AuthService struct {
	store *store.Queries
	jwt   *auth.JWTMiddleware
}

func NewAuthService(s *store.Queries, jwt *auth.JWTMiddleware) *AuthService {
	return &AuthService{store: s, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	dbu, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbu.PasswordHash), []byte(password)); err != nil {
		return "", err
	}
	return s.jwt.GenerateToken(dbu.ID, dbu.Username)
}
