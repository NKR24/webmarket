package product

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateProduct(p *Product) (uuid.UUID, error) {
	p.ID = uuid.New()
	const query = `INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, p.Name, p.Price, p.Stock).Scan(&p.ID)
	if err != nil {
		return p.ID, err
	}
	return p.ID, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]Product, error) {
	var products []Product
	const query = "SELECT id, name, price, stock FROM products"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		product := new(Product)
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, nil
}

func (r *Repository) GetById(id uuid.UUID) (*Product, error) {
	product := new(Product)
	const query = `SELECT id, name, price, stock FROM products WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *Repository) Update(id uuid.UUID, p *Product) error {
	const query = `UPDATE products SET name = $1, price = $2, stock = $3  WHERE id = $4`
	_, err := r.db.Exec(query, p.Name, p.Price, p.Stock, id)
	return err
}

func (r *Repository) Delete(id uuid.UUID) error {
	const query = `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
