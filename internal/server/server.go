package server

import (
	"car_factory/internal/handlers"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartServer(db *pgxpool.Pool) {
	router := mux.NewRouter()

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/index.html")
	})
	api := router.PathPrefix("/api").Subrouter()

	api.HandleFunc("/suppliers", handlers.GetSuppliersID(db)).Methods("GET")
	api.HandleFunc("/suppliers/report", handlers.SupplierReportHandler(db)).Methods("GET")
	api.HandleFunc("/cars", handlers.GetCarsHandler(db)).Methods("GET")
	api.HandleFunc("/parts", handlers.GetAllPartsHandlers(db)).Methods("GET")
	api.HandleFunc("/parts/stock", handlers.GetPartStorageHandler(db)).Methods("GET")
	api.HandleFunc("/equipment/usage", handlers.EquipmentUsageHandler(db)).Methods("GET")
	api.HandleFunc("/production", handlers.ProductionReportHandler(db)).Methods("GET")
	api.HandleFunc("/production/chart", handlers.ProductionChartHandler(db)).Methods("GET")

	server := &http.Server{
		Addr:         ":9091",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Сервер запущен на :9091")
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Ошибка: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	gracefulShutdown(quit, server)
}
