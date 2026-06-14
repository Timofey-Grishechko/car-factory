package models

import (
	"time"
)

type ProductionReport struct {
	ProductionDate time.Time `json:"production_date"`
	Planned        int       `json:"planned"`
	Actual         int       `json:"actual"`
	PercentDone    float64   `json:"percent_done"`
}

type EquipmentStat struct {
	MachineNumber int    `json:"machine_number"`
	Name          string `json:"name"`
	TotalHours    int    `json:"total_hours"`
}

type PartStock struct {
	PartCode int    `json:"part_code"`
	Name     string `json:"name"`
	Supplied int    `json:"supplied"`
	Consumed int    `json:"consumed"`
	Balance  int    `json:"balance"`
}

type PartItem struct {
	PartCode int    `json:"part_code"`
	Name     string `json:"name"`
}

type CarInfo struct {
	VinNumber      string    `json:"vin_number"`
	Color          string    `json:"color"`
	ProductionDate time.Time `json:"production_date"`
	Model          string    `json:"model"`
}

type SupplierPartDetail struct {
	PartCode      int     `json:"part_code"`
	PartName      string  `json:"part_name"`
	TotalQuantity int     `json:"total_quantity"`
	TotalAmount   float64 `json:"total_amount"`
}

type MonthlySupplySummary struct {
	YearMonth     string  `json:"year_month"`
	TotalAmount   float64 `json:"total_amount"`
	TotalQuantity int     `json:"total_quantity"`
	SupplyCount   int     `json:"supply_count"`
}

type SupplierReport struct {
	SupplierID   int                    `json:"supplier_ID"`
	SupplierName string                 `json:"supplier_Name"`
	Parts        []SupplierPartDetail   `json:"parts"`
	Monthly      []MonthlySupplySummary `json:"monthly"`
}

type Supplier struct {
	SupplierID   int    `json:"supplier_ID"`
	SupplierName string `json:"supplier_Name"`
}

type ProductionResponse struct {
	Data   []ProductionReport `json:"data"`
	Total  int                `json:"total"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

func NewProductionResponse(data []ProductionReport, total, limit, offset int) ProductionResponse {
	return ProductionResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}

type CarResponse struct {
	Data   []CarInfo `json:"data"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}

func NewCarResponse(data []CarInfo, total, limit, offset int) CarResponse {
	return CarResponse{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}
}
