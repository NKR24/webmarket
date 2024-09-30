package main

import (
	"database/sql"
	"log"
	"shop/config"
	"shop/internal/product"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := sql.Open("postgres",
		"host="+cfg.DBHost+" port="+cfg.DBPort+" user="+cfg.DBUser+" password="+cfg.DBPassword+" dbname="+cfg.DBName+" sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	kafkaWriter := product.NewKafkaWriter(cfg.KafkaBrokers, cfg.KafkaTopic, 1)
	// if err != nil {
	// 	panic(err)
	// }

	repo := product.NewRepository(db)

	handler := product.NewHandler(repo, kafkaWriter)

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.POST("/product", handler.CreateNewProduct)
	e.GET("/products", handler.GetAllProducts)
	e.GET("/product/:id", handler.GetProductById)
	e.PUT("product/update/:id", handler.UpdateProductById)
	e.DELETE("product/delete/:id", handler.DeleteProductById)
	e.Logger.Fatal(e.Start(":9090"))
}
