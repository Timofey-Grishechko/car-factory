package handlers

import (
	"car_factory/internal/models"
	"car_factory/internal/repository"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCarsHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		color := r.URL.Query().Get("color")
		model := r.URL.Query().Get("model")
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		limit, offset := getPaginationParams(r)

		var from, to time.Time
		var err error
		if dateFrom != "" {
			from, err = time.Parse("2006-01-02", dateFrom)
			if err != nil {
				http.Error(w, "invalid date_from, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}
		if dateTo != "" {
			to, err = time.Parse("2006-01-02", dateTo)
			if err != nil {
				http.Error(w, "invalid date_to, use YYYY-MM-DD", http.StatusBadRequest)
				return
			}
		}

		cars, total, err := repository.GetCar(r.Context(), db, model, color, from, to, limit, offset)
		if err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		response := models.NewCarResponse(cars, total, limit, offset)
		responseJSON(w, http.StatusOK, response)
	}
}
