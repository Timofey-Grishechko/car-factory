package handlers

import (
	"car_factory/internal/repository"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetSuppliersID(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		suppliers, err := repository.GetSuppliers(r.Context(), db)
		if err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		responseJSON(w, http.StatusOK, suppliers)
	}
}

func SupplierReportHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supplierIDStr := r.URL.Query().Get("supplier_id")
		var supplierID int
		var err error
		if supplierIDStr != "" {
			supplierID, err = strconv.Atoi(supplierIDStr)
			if err != nil || supplierID <= 0 {
				http.Error(w, "invalid supplier_id, must be positive integer", http.StatusBadRequest)
				return
			}
		}

		report, err := repository.GetSupplierReport(r.Context(), db, supplierID)
		if err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if supplierID > 0 && len(report) == 0 {
			http.Error(w, "supplier wasn't found", http.StatusNotFound)
			return
		}

		responseJSON(w, http.StatusOK, report)
	}
}
