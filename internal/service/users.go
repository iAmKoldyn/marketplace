package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/iAmKoldyn/marketplace/internal/domain"
	"github.com/iAmKoldyn/marketplace/internal/store/sqlc"
)

type UserService struct {
	store *sqlc.Queries
}

func NewUserService(s *sqlc.Queries) *UserService {
	return &UserService{store: s}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	dbu, err := s.store.CreateUser(ctx, sqlc.CreateUserParams{
		Username:     username,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        int64(dbu.ID),
		Username:  dbu.Username,
		CreatedAt: dbu.CreatedAt,
	}, nil
}
