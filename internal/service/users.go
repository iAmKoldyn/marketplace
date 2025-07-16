package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/yourusername/marketplace/internal/domain"
	"github.com/yourusername/marketplace/internal/store"
)

type UserService struct {
	store *store.Queries
}

func NewUserService(s *store.Queries) *UserService {
	return &UserService{store: s}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*domain.User, error) {
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	dbu, err := s.store.CreateUser(ctx, store.CreateUserParams{
		Username:     username,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        dbu.ID,
		Username:  dbu.Username,
		CreatedAt: dbu.CreatedAt,
	}, nil
}
