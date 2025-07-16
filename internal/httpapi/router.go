package httpapi

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/yourusername/marketplace/internal/auth"
	"github.com/yourusername/marketplace/internal/config"
	"github.com/yourusername/marketplace/internal/middleware"
	"github.com/yourusername/marketplace/internal/service"
	"github.com/yourusername/marketplace/internal/store"
)

func NewRouter(db *store.Queries, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// built-in
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer, middleware.Timeout(60*time.Second))

	// our custom
	r.Use(middleware.CORS())
	r.Use(middleware.RequireJSON)
	r.Use(middleware.RateLimit)
	r.Use(middleware.Metrics)

	// metrics endpoint
	r.Handle("/metrics", middleware.MetricsHandler())

	// auth & routes…
	jwtm := auth.NewJWTMiddleware(cfg.JWTSecret, cfg.JWTExpiry)
	us := service.NewUserService(db)
	as := service.NewAuthService(db, jwtm)
	ads := service.NewAdsService(db)

	r.Post("/users/register", handlersUsers(us))
	r.Post("/auth/login", handlersAuth(as))

	r.Group(func(r chi.Router) {
		r.Use(jwtm.Verifier, jwtm.Authenticator)
		r.Post("/ads", handlersCreateAd(ads))
		r.Get("/ads", handlersListAds(ads))
	})
	return r
}
