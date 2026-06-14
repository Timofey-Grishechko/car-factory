package handlers

import (
	"car_factory/internal/repository"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EquipmentUsageHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		machineType := r.URL.Query().Get("machine_type")
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")

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

		stats, err := repository.EquipmentUsage(r.Context(), db, machineType, from, to)
		if err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		responseJSON(w, http.StatusOK, stats)
	}
}
