package service

import (
	"context"
	"math"

	"github.com/yourusername/marketplace/internal/domain"
	"github.com/yourusername/marketplace/internal/store"
)

type AdsService struct {
	store *store.Queries
}

func NewAdsService(s *store.Queries) *AdsService {
	return &AdsService{store: s}
}

func (s *AdsService) Create(ctx context.Context, authorID int64, title, text, imageURL string, price float64) (*domain.Ad, error) {
	dbad, err := s.store.CreateAd(ctx, store.CreateAdParams{
		AuthorID: authorID,
		Title:    title,
		Text:     text,
		ImageURL: imageURL,
		Price:    price,
	})
	if err != nil {
		return nil, err
	}
	return &domain.Ad{
		ID:             dbad.ID,
		AuthorID:       dbad.AuthorID,
		AuthorUsername: "", // filled in list
		Title:          dbad.Title,
		Text:           dbad.Text,
		ImageURL:       dbad.ImageURL,
		Price:          float64(dbad.Price),
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
	var dbads []store.ListAdsByDateDescRow
	// dispatch to correct query
	switch sortBy + "_" + order {
	case "date_desc":
		dbads, _ = s.store.ListAdsByDateDesc(ctx, minPrice, maxPrice, int32(pageSize), int32(offset))
	case "date_asc":
		dbads, _ = s.store.ListAdsByDateAsc(ctx, minPrice, maxPrice, int32(pageSize), int32(offset))
	case "price_desc":
		dbads, _ = s.store.ListAdsByPriceDesc(ctx, minPrice, maxPrice, int32(pageSize), int32(offset))
	default:
		dbads, _ = s.store.ListAdsByPriceAsc(ctx, minPrice, maxPrice, int32(pageSize), int32(offset))
	}

	out := make([]domain.Ad, len(dbads))
	for i, d := range dbads {
		out[i] = domain.Ad{
			ID:             d.ID,
			AuthorID:       d.AuthorID,
			AuthorUsername: d.AuthorUsername,
			Title:          d.Title,
			Text:           d.Text,
			ImageURL:       d.ImageURL,
			Price:          float64(d.Price),
			CreatedAt:      d.CreatedAt,
			IsMine:         d.AuthorID == currentUserID,
		}
	}
	return out, nil
}
