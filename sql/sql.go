package sql

var PartQuery string = `
	WITH supplied AS (
		SELECT part_code, SUM(quantity) AS total_supplied
		FROM supply_content
		GROUP BY part_code
	),
	used AS (
		SELECT mp.part_code, SUM(c.car_count * mp.quantity) AS total_used
		FROM model_part mp
		JOIN (
			SELECT model_code, COUNT(*) AS car_count
			FROM car
			GROUP BY model_code
		) c ON mp.model_code = c.model_code
		GROUP BY mp.part_code
	)
	SELECT
		p.part_code,
		p.name,
		COALESCE(s.total_supplied, 0) AS supplied,
		COALESCE(u.total_used, 0) AS consumed,
		COALESCE(s.total_supplied, 0) - COALESCE(u.total_used, 0) AS balance
	FROM part p
	LEFT JOIN supplied s ON p.part_code = s.part_code
	LEFT JOIN used u ON p.part_code = u.part_code
`
var SupplyQuery string = `
	SELECT
		sc.part_code,
		p.name,
		SUM(sc.quantity) AS total_qty,
		SUM(sc.quantity * sc.price_per_unit) AS total_amount
	FROM supply_content sc
	JOIN supply s ON sc.supply_number = s.supply_number
	JOIN part p ON sc.part_code = p.part_code
	WHERE s.supplier_id = $1
	GROUP BY sc.part_code, p.name
	ORDER BY total_amount DESC
`
var MonthlyQuery string = `
	SELECT
		TO_CHAR(s.supply_date, 'YYYY-MM') AS year_month,
		COALESCE(SUM(sc.quantity * sc.price_per_unit), 0) AS total_amount,
		COALESCE(SUM(sc.quantity), 0) AS total_quantity,
		COUNT(DISTINCT s.supply_number) AS supply_count
	FROM supply s
	JOIN supply_content sc ON s.supply_number = sc.supply_number
	WHERE s.supplier_id = $1
	GROUP BY year_month
	ORDER BY year_month
`
