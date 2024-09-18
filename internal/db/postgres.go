package db

import (
	"database/sql"
	"market/internal/kafka"
	"market/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DataBase struct {
	Conn *sql.DB
}

func Connect(dsn string) (*DataBase, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DataBase{Conn: db}, nil
}

func (db *DataBase) InsertProduct(p *models.Product) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO products (name, description, stock, quantity, price) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.Conn.QueryRow(query, p.Name, p.Description, p.Stock, p.Quantity, p.Price).Scan(&id)
	if err != nil {
		return "Insertion error", nil
	}
	return id, nil
}

func (db *DataBase) SelectProduct(id string, kp *kafka.KafkaProducer) (*models.Product, error) {
	query := `SELECT * FROM products WHERE id = $1`
	row := db.Conn.QueryRow(query, id)
	p := new(models.Product)
	err := row.Scan(p.ID, p.Name, p.Description, p.Stock, p.Quantity, p.Price)
	if err != nil {
		return nil, err
	}
	kp.SendMessage("ProductSelected: ", p.ID)
	return p, nil
}

func (db *DataBase) UpdateProduct(p *models.Product, kp *kafka.KafkaProducer) error {
	query := `UPDATE products SET name = $1, description = $2, stock = $3, quantity = $4, price = $5 WHERE id = $6`
	_, err := db.Conn.Exec(query, p.Name, p.Description, p.Stock, p.Quantity, p.Price)
	kp.SendMessage("ProductUpdated", p)
	return err
}

func (db *DataBase) DeleteProduct(id string, kp *kafka.KafkaProducer) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := db.Conn.Exec(query, id)
	kp.SendMessage("ProductDeleted", id)
	return err
}

// Пример отправки события после создания продукта
func (db *DataBase) CreateProductWithKafka(p *models.Product, kp *kafka.KafkaProducer) (string, error) {
	id, err := db.InsertProduct(p)
	if err != nil {
		return "", err
	}
	// Отправляем сообщение в Kafka
	kp.SendMessage("ProductCreated", p)
	return id, nil
}
