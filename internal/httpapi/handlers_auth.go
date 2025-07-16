package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/marketplace/internal/service"
)

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handlersAuth(as *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		token, err := as.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		JSON(w, http.StatusOK, map[string]string{"token": token})
	}
}
