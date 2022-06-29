package model

type Order struct {
	Number string
	Status string
}

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
