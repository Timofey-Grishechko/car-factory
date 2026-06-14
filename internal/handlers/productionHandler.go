package handlers

import (
	"car_factory/internal/models"
	"car_factory/internal/repository"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductionReportHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := getPaginationParams(r)
		reports, total, err := repository.GetProductionReport(r.Context(), pool, limit, offset)
		if err != nil {
			http.Error(w, "failed to fetch production report: "+err.Error(), http.StatusInternalServerError)
			return
		}

		productionResponse := models.NewProductionResponse(reports, total, limit, offset)

		responseJSON(w, http.StatusOK, productionResponse)
	}

}

func ProductionChartHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		report, err := repository.GetAllProduction(r.Context(), pool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseJSON(w, http.StatusOK, report)
	}
}
