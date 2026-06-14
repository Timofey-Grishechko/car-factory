package repository

import (
	"car_factory/internal/models"
	"car_factory/sql"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllParts(ctx context.Context, db *pgxpool.Pool) ([]models.PartItem, error) {
	rows, err := db.Query(ctx, "SELECT part_code, name FROM part ORDER BY part_code")
	if err != nil {
		log.Printf("GetAllParts при запросе: %v", err)
		return nil, err
	}
	defer rows.Close()

	var parts []models.PartItem
	for rows.Next() {
		var p models.PartItem
		if err := rows.Scan(&p.PartCode, &p.Name); err != nil {
			log.Printf("GetAllParts при обработке: %v", err)
			return nil, err
		}

		parts = append(parts, p)
	}

	return parts, nil
}

func GetPartStorage(ctx context.Context, db *pgxpool.Pool, partCode int) ([]models.PartStock, error) {
	query := sql.PartQuery

	args := []interface{}{}
	if partCode > 0 {
		query += " WHERE p.part_code = $1"
		args = append(args, partCode)
	}
	query += " ORDER BY p.part_code"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("GetPartStorage при запросе: %v", err)
		return nil, err
	}
	defer rows.Close()

	var stocks []models.PartStock
	for rows.Next() {
		var s models.PartStock
		if err := rows.Scan(&s.PartCode, &s.Name, &s.Supplied, &s.Consumed, &s.Balance); err != nil {
			log.Printf("GetPartStorage при обработке: %v", err)
			return nil, err
		}

		stocks = append(stocks, s)
	}

	return stocks, nil
}
