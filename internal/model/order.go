package model

import "time"

type AccrualOrder struct {
	UserID     int       `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type WithdrawOrder struct {
	UserID      int       `json:"-"`
	Order       string    `json:"order"`
	Sum         float32   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}
