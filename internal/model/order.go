package model

import "time"

type Order struct {
	Number     string    `json:"number"`
	UserID     int       `json:"user_id"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
