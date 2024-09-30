package product

import (
	"context"

	"github.com/google/uuid"
)

type Repositer interface {
	CreateProduct(p *Product) (uuid.UUID, error)
	GetAll(ctx context.Context) ([]Product, error)
	GetById(id uuid.UUID) (*Product, error)
	Update(id uuid.UUID, p *Product) error
	Delete(id uuid.UUID) error
}
