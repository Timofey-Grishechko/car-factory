package handlers

import (
	"car_factory/internal/repository"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllPartsHandlers(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts, err := repository.GetAllParts(r.Context(), db)
		if err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		responseJSON(w, http.StatusOK, parts)
	}
}

func GetPartStorageHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		partCodeStr := r.URL.Query().Get("part_code")

		var partCode int
		var err error
		if partCodeStr != "" {
			partCode, err = strconv.Atoi(partCodeStr)
			if err != nil {
				http.Error(w, "invalid part_code, must be positive integer", http.StatusBadRequest)
				return
			}
		}

		stocks, err := repository.GetPartStorage(r.Context(), db, partCode)
		if err != nil {
			http.Error(w, "part not found", http.StatusNotFound)
			return
		}

		responseJSON(w, http.StatusOK, stocks)
	}
}
