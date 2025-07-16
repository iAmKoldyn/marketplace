package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/marketplace/internal/service"
	"github.com/yourusername/marketplace/internal/validation"
)

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handlersUsers(us *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if err := validation.ValidateUsername(req.Username); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := validation.ValidatePassword(req.Password); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := us.Register(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		JSON(w, http.StatusCreated, user)
	}
}
