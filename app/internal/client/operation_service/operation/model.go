package operation

import "time"

type Operation struct {
	UUID         string    `json:"uuid"`
	CategoryUUID string    `json:"category_uuid"`
	MoneySum     float64   `json:"money_sum"`
	Description  string    `json:"description"`
	DateTime     time.Time `json:"date_time"`
}

type CreateOperationDTO struct {
	CategoryUUID string  `json:"category_uuid"`
	MoneySum     float64 `json:"money_sum"`
	Description  string  `json:"description"`
}

type UpdateOperationDTO struct {
	CategoryUUID string  `json:"category_uuid"`
	MoneySum     float64 `json:"money_sum"`
	Description  string  `json:"description"`
}
