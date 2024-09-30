package product

import "github.com/google/uuid"

type Product struct {
	ID    uuid.UUID `json:"uuid"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
	Stock int       `json:"stock"`
}
