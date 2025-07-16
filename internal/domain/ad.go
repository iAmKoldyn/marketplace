package domain

import "time"

type Ad struct {
	ID             int64     `json:"id"`
	AuthorID       int64     `json:"-"`
	AuthorUsername string    `json:"author_username"`
	Title          string    `json:"title"`
	Text           string    `json:"text"`
	ImageURL       string    `json:"image_url"`
	Price          float64   `json:"price"`
	CreatedAt      time.Time `json:"created_at"`
	IsMine         bool      `json:"is_mine,omitempty"`
}
