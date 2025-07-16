package httpapi

import (
	"time"

	"github.com/go-chi/chi/v5"
	chiMw "github.com/go-chi/chi/v5/middleware"

	"github.com/iAmKoldyn/marketplace/internal/auth"
	"github.com/iAmKoldyn/marketplace/internal/config"
	imw "github.com/iAmKoldyn/marketplace/internal/middleware"
	"github.com/iAmKoldyn/marketplace/internal/service"
	"github.com/iAmKoldyn/marketplace/internal/store/sqlc"
)

func NewRouter(db *sqlc.Queries, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// built-in
	r.Use(chiMw.RequestID, chiMw.Logger, chiMw.Recoverer, chiMw.Timeout(60*time.Second))

	// our custom
	r.Use(imw.CORS())
	r.Use(imw.RequireJSON)
	r.Use(imw.RateLimit)
	r.Use(imw.Metrics)

	// metrics endpoint
	r.Handle("/metrics", imw.MetricsHandler())

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
