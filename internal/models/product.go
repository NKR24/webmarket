package models

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stock       bool      `json:"stock"`
	Quantity    int64     `json:"quantity"`
	Price       float64   `json:"price"`
}
