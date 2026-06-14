package repository

import (
	"car_factory/internal/models"
	"car_factory/sql"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetSuppliers(ctx context.Context, db *pgxpool.Pool) ([]models.Supplier, error) {
	query := `
			SELECT supplier_id, supplier_name
			FROM supplier
		`

	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Printf("GetSuppliers при запросе: %v", err)
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(&s.SupplierID, &s.SupplierName); err != nil {
			log.Printf("GetSuppliers при обработке: %v", err)
			return nil, err
		}

		suppliers = append(suppliers, s)
	}

	return suppliers, nil

}

func GetSupplierReport(ctx context.Context, db *pgxpool.Pool, supplierID int) ([]models.SupplierReport, error) {
	suppliersQuery := `SELECT supplier_id, supplier_name FROM supplier`
	args := []interface{}{}
	if supplierID > 0 {
		suppliersQuery += " WHERE supplier_id = $1"
		args = append(args, supplierID)
	}
	suppliersQuery += " ORDER BY supplier_id"

	rows, err := db.Query(ctx, suppliersQuery, args...)
	if err != nil {
		log.Printf("GetSupplierReport при запроса: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reports []models.SupplierReport
	for rows.Next() {
		var s models.SupplierReport
		if err := rows.Scan(&s.SupplierID, &s.SupplierName); err != nil {
			log.Printf("GetSupplierReport при обработке: %v", err)
			return nil, err
		}

		partRows, err := db.Query(ctx, sql.SupplyQuery, s.SupplierID)
		if err != nil {
			log.Printf("GetSupplierReport при запросе partrows: %v", err)
			return nil, err
		}

		var parts []models.SupplierPartDetail
		for partRows.Next() {
			var p models.SupplierPartDetail
			if err := partRows.Scan(&p.PartCode, &p.PartName, &p.TotalQuantity, &p.TotalAmount); err != nil {
				partRows.Close()
				log.Printf("GetSupplierReport при обработке partrows: %v", err)
				return nil, err
			}

			parts = append(parts, p)

		}

		partRows.Close()
		s.Parts = parts

		monthRows, err := db.Query(ctx, sql.MonthlyQuery, s.SupplierID)
		if err != nil {
			log.Printf("GetSupplierReport при запросе monthrows: %v", err)
			return nil, err
		}
		var monthly []models.MonthlySupplySummary
		for monthRows.Next() {
			var m models.MonthlySupplySummary
			if err := monthRows.Scan(&m.YearMonth, &m.TotalAmount, &m.TotalQuantity, &m.SupplyCount); err != nil {
				monthRows.Close()
				log.Printf("GetSupplierReport при обработке monthrows: %v", err)
				return nil, err
			}
			monthly = append(monthly, m)
		}

		monthRows.Close()
		s.Monthly = monthly

		reports = append(reports, s)
	}
	return reports, nil
}
