package stats_service

import "time"

type Operation struct {
	UUID         string    `json:"uuid"`
	CategoryUUID string    `json:"category_uuid"`
	Description  string    `json:"description"`
	MoneySum     float64   `json:"money_sum"`
	DateTime     time.Time `json:"date_time"`
}

type Report struct {
	TotalMoneySum float64     `json:"total_money_sum"`
	Operations    []Operation `json:"operations"`
}
