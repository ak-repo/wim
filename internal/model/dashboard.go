package model

type TotalCount struct {
	TotalProducts   int `json:"total_products"`
	TotalUsers      int `json:"total_users"`
	TotalWarehouses int `json:"total_warehouses"`
	TotalLocations  int `json:"total_locations"`
}
