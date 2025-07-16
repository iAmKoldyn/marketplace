package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"os"

	"github.com/hibiken/asynq"
	"github.com/yourusername/marketplace/internal/auth"
	"github.com/yourusername/marketplace/internal/service"
	"github.com/yourusername/marketplace/internal/validation"
	"github.com/yourusername/marketplace/internal/worker"
)

type createAdReq struct {
	Title    string  `json:"title"`
	Text     string  `json:"text"`
	ImageURL string  `json:"image_url"`
	Price    float64 `json:"price"`
}

func handlersCreateAd(as *service.AdsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createAdReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		if err := validation.ValidateAdTitle(req.Title); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := validation.ValidateAdText(req.Text); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := validation.ValidateImageURL(req.ImageURL); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := validation.ValidatePrice(req.Price); err != nil {
			Error(w, http.StatusBadRequest, err.Error())
			return
		}
		claims := auth.FromContext(r.Context())
		ad, err := as.Create(r.Context(), claims.UserID, req.Title, req.Text, req.ImageURL, req.Price)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		JSON(w, http.StatusCreated, ad)

		task, err := worker.NewGenerateThumbnailTask(ad.ID, ad.ImageURL)
		if err == nil {
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR")})
			defer client.Close()
			client.Enqueue(task, asynq.Queue("default"))
		}
	}
}

func handlersListAds(as *service.AdsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		minP, _ := strconv.ParseFloat(q.Get("min_price"), 64)
		maxP, _ := strconv.ParseFloat(q.Get("max_price"), 64)
		sortBy := q.Get("sort_by") // "date" or "price"
		order := q.Get("order")    // "asc" or "desc"
		if sortBy != "price" {
			sortBy = "date"
		}
		if order != "asc" {
			order = "desc"
		}
		page, _ := strconv.Atoi(q.Get("page"))
		if page < 1 {
			page = 1
		}
		per, _ := strconv.Atoi(q.Get("per_page"))
		if per < 1 || per > 100 {
			per = 20
		}

		var userID int64
		if claims := auth.FromContext(r.Context()); claims != nil {
			userID = claims.UserID
		}

		ads, err := as.List(r.Context(), minP, maxP, sortBy, order, page, per, userID)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		JSON(w, http.StatusOK, ads)
	}
}
