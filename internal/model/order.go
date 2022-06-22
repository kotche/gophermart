package model

type Order struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	User   string
	Status string
}
