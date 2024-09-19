package product

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePruduct(product Product) error {
	query := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, product.Name, product.Price).Scan(&product.ID)
}

func (r *Repository) GetAll() ([]Product, error) {
	rows, err := r.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 3, 6)
	for rows.Next() {
		product := new(Product)
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, nil
}

func (r *Repository) GetById(id int) (*Product, error) {
	product := new(Product)
	err := r.db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *Repository) Update(product *Product) error {
	_, err := r.db.Exec("UPDATE products SET name = $1, price = $2 WHERE id = $3", product.Name, product.Price, product.ID)
	return err
}

func (r *Repository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
