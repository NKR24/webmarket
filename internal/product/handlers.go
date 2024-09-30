package product

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Handler struct {
	repo Repositer
	kw   *KafkaWriter
}

func NewHandler(repo Repositer, kw *KafkaWriter) *Handler {
	return &Handler{
		repo: repo,
		kw:   kw,
	}
}

func (h *Handler) CreateNewProduct(c echo.Context) error {
	product := new(Product)
	if err := c.Bind(&product); err != nil {
		return err
	}
	fmt.Printf("Input data: %+v\n", product)

	if _, err := h.repo.CreateProduct(product); err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	if err := h.kw.SendMessage(product); err != nil {
		log.Printf("Error to send message to kafka: %v", err)
		return err
	}

	return c.JSON(
		http.StatusOK,
		product.ID,
	)
}

func (h *Handler) GetAllProducts(c echo.Context) error {
	products, err := h.repo.GetAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, products)
}

func (h *Handler) GetProductById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	log.Printf("id %s", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	product, err := h.repo.GetById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, product)
}

func (h *Handler) UpdateProductById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	product := new(Product)
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	product.ID = id

	err = h.repo.Update(id, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, product)
}

func (h *Handler) DeleteProductById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = h.repo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Sucessfuly Deleted!")
}
