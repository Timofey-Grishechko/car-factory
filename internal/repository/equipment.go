package repository

import (
	"car_factory/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func EquipmentUsage(
	ctx context.Context,
	db *pgxpool.Pool,
	machineType string,
	dateFrom,
	dateTo time.Time,
) ([]models.EquipmentStat, error) {

	query := `
		SELECT e.machine_number, e.name, COALESCE(SUM(mu.hours_worked), 0) as total_hours
		FROM equipment e
		LEFT JOIN machine_usage mu ON e.machine_number = mu.machine_number
		WHERE 1=1
	`
	args := []interface{}{}
	argID := 1

	if machineType != "" {
		query += fmt.Sprintf(" AND e.machine_type = $%d", argID)
		args = append(args, machineType)
		argID++
	}
	if !dateFrom.IsZero() {
		query += fmt.Sprintf(" AND mu.production_date >= $%d", argID)
		args = append(args, dateFrom)
		argID++
	}
	if !dateTo.IsZero() {
		query += fmt.Sprintf(" AND mu.production_date <= $%d", argID)
		args = append(args, dateTo)
		argID++
	}

	query += " GROUP BY e.machine_number, e.name ORDER BY e.machine_number"

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		log.Printf("EquipmentUsage при запросе: %v", err)
		return nil, err
	}
	defer rows.Close()

	var stats []models.EquipmentStat
	for rows.Next() {
		var s models.EquipmentStat
		if err := rows.Scan(&s.MachineNumber, &s.Name, &s.TotalHours); err != nil {
			log.Printf("EquipmentUsage при обработке: %v", err)
			return nil, err
		}
		stats = append(stats, s)
	}

	return stats, nil
}
