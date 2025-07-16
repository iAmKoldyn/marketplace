package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userCtxKey = contextKey("user")

type JWTMiddleware struct {
	secret string
	expiry time.Duration
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTMiddleware(secret string, minutes int) *JWTMiddleware {
	return &JWTMiddleware{
		secret: secret,
		expiry: time.Duration(minutes) * time.Minute,
	}
}

func (j *JWTMiddleware) GenerateToken(userID int64, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(j.secret))
}

// Verifier parses token if present and stores claims in context.
func (j *JWTMiddleware) Verifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := r.Header.Get("Authorization")
		if hdr == "" {
			next.ServeHTTP(w, r)
			return
		}
		parts := strings.SplitN(hdr, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusBadRequest)
			return
		}
		tokenStr := parts[1]
		tok, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(j.secret), nil
		})
		if err == nil {
			if claims, ok := tok.Claims.(*Claims); ok && tok.Valid {
				ctx := context.WithValue(r.Context(), userCtxKey, claims)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Authenticator requires that a valid token be present.
func (j *JWTMiddleware) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(userCtxKey).(*Claims)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// pass through
		ctx := context.WithValue(r.Context(), userCtxKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext retrieves the JWT claims
func FromContext(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(userCtxKey).(*Claims); ok {
		return claims
	}
	return nil
}
