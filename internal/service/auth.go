package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/iAmKoldyn/marketplace/internal/auth"
	"github.com/iAmKoldyn/marketplace/internal/store/sqlc"
)

type AuthService struct {
	store *sqlc.Queries
	jwt   *auth.JWTMiddleware
}

func NewAuthService(s *sqlc.Queries, jwtm *auth.JWTMiddleware) *AuthService {
	return &AuthService{store: s, jwt: jwtm}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	dbu, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbu.PasswordHash), []byte(password)); err != nil {
		return "", err
	}
	return s.jwt.GenerateToken(int64(dbu.ID), dbu.Username)
}
