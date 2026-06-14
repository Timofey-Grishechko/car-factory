package repository

import (
	"car_factory/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCar(ctx context.Context, db *pgxpool.Pool, model, color string, from, to time.Time, limit, offset int) ([]models.CarInfo, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM car c
		LEFT JOIN model m ON c.model_code = m.model_code
		WHERE 1=1
	`

	dataQuery := `
		SELECT c.vin_number, c.color, c.production_date, m.name
		FROM car c
		LEFT JOIN model m ON c.model_code = m.model_code
		WHERE 1=1
	`
	var args []interface{}
	argID := 1

	addCondition := func(condition string, value interface{}) {
		countQuery += condition
		dataQuery += condition
		args = append(args, value)
		argID += 1
	}

	if model != "" {
		addCondition(fmt.Sprintf(" AND m.name ILIKE $%d", argID), "%"+model+"%")
	}
	if color != "" {
		addCondition(fmt.Sprintf(" AND c.color ILIKE $%d", argID), "%"+color+"%")
	}
	if !from.IsZero() {
		addCondition(fmt.Sprintf(" AND c.production_date >= $%d", argID), from)
	}
	if !to.IsZero() {
		addCondition(fmt.Sprintf(" AND c.production_date <= $%d", argID), to)
	}

	var total int
	err := db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("GetCar при подсчете: %v", err)
		return nil, 0, err
	}

	dataQuery += " ORDER BY c.production_date DESC LIMIT $%d OFFSET $%d"
	dataQuery = fmt.Sprintf(dataQuery, argID, argID+1)
	args = append(args, limit, offset)

	rows, err := db.Query(ctx, dataQuery, args...)
	if err != nil {
		log.Printf("GetCar при запросе: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var cars []models.CarInfo
	for rows.Next() {
		var c models.CarInfo
		if err := rows.Scan(&c.VinNumber, &c.Color, &c.ProductionDate, &c.Model); err != nil {
			log.Printf("GetCar при обработке: %v", err)
			return nil, 0, err
		}
		cars = append(cars, c)
	}
	return cars, total, nil
}
