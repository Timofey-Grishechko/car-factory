package repository

import (
	"car_factory/internal/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProductionReport(ctx context.Context, db *pgxpool.Pool, limit, offset int) ([]models.ProductionReport, int, error) {
	var total int
	err := db.QueryRow(ctx, "SELECT COUNT(*) FROM production").Scan(&total)
	if err != nil {
		log.Printf("GetProductionReport при запросе количества строк: %v", err)
		return nil, 0, err
	}

	rows, err := db.Query(ctx, `
		SELECT production_date, planned_quantity, actual_quantity,
			COALESCE(ROUND(actual_quantity::numeric / NULLIF(planned_quantity, 0) * 100, 2), 0) as percent
		FROM production
		ORDER BY production_date
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		log.Printf("GetProductionReport при запросе: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var reports []models.ProductionReport
	for rows.Next() {
		var r models.ProductionReport
		if err := rows.Scan(&r.ProductionDate, &r.Planned, &r.Actual, &r.PercentDone); err != nil {
			log.Printf("GetProductionReport при обработке запроса: %v", err)
			return nil, 0, err
		}
		reports = append(reports, r)
	}

	return reports, total, nil
}

func GetAllProduction(ctx context.Context, db *pgxpool.Pool) ([]models.ProductionReport, error) {
	query := `
		SELECT production_date, planned_quantity, actual_quantity,
		COALESCE(ROUND(actual_quantity::numeric / NULLIF(planned_quantity, 0)*100, 2), 0) as percent
		FROM production ORDER BY production_date
		`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.ProductionReport
	for rows.Next() {
		var r models.ProductionReport
		if err := rows.Scan(&r.ProductionDate, &r.Planned, &r.Actual, &r.PercentDone); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}
