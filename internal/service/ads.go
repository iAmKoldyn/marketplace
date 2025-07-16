package service

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/iAmKoldyn/marketplace/internal/domain"
	"github.com/iAmKoldyn/marketplace/internal/store/sqlc"
)

type AdsService struct {
	store *sqlc.Queries
}

func NewAdsService(s *sqlc.Queries) *AdsService {
	return &AdsService{store: s}
}

func (s *AdsService) Create(ctx context.Context, authorID int64, title, text, imageURL string, price float64) (*domain.Ad, error) {
	dbad, err := s.store.CreateAd(ctx, sqlc.CreateAdParams{
		AuthorID: int32(authorID),
		Title:    title,
		Text:     text,
		ImageUrl: imageURL,
		Price:    fmt.Sprintf("%f", price),
	})
	if err != nil {
		return nil, err
	}
	priceFloat, _ := strconv.ParseFloat(dbad.Price, 64)
	return &domain.Ad{
		ID:             int64(dbad.ID),
		AuthorID:       int64(dbad.AuthorID),
		AuthorUsername: "", // filled in list
		Title:          dbad.Title,
		Text:           dbad.Text,
		ImageURL:       dbad.ImageUrl,
		Price:          priceFloat,
		CreatedAt:      dbad.CreatedAt,
	}, nil
}

// List with pagination, sorting and price filtering.
func (s *AdsService) List(ctx context.Context,
	minPrice, maxPrice float64,
	sortBy, order string,
	page, pageSize int,
	currentUserID int64,
) ([]domain.Ad, error) {
	if maxPrice < minPrice {
		maxPrice = math.Max(minPrice, maxPrice)
	}
	offset := (page - 1) * pageSize
	params := sqlc.ListAdsByDateDescParams{
		Price:   fmt.Sprintf("%f", minPrice),
		Price_2: fmt.Sprintf("%f", maxPrice),
		Limit:   int32(pageSize),
		Offset:  int32(offset),
	}
	dbads, err := s.store.ListAdsByDateDesc(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Ad, len(dbads))
	for i, d := range dbads {
		priceFloat, _ := strconv.ParseFloat(d.Price, 64)
		out[i] = domain.Ad{
			ID:             int64(d.ID),
			AuthorID:       int64(d.AuthorID),
			AuthorUsername: d.AuthorUsername,
			Title:          d.Title,
			Text:           d.Text,
			ImageURL:       d.ImageUrl,
			Price:          priceFloat,
			CreatedAt:      d.CreatedAt,
			IsMine:         int64(d.AuthorID) == currentUserID,
		}
	}
	return out, nil
}
