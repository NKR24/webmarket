package product

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type Handler struct {
	repo *Repository
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

func (h *Handler) CreateNewProduct(c echo.Context) error {
	product := new(Product)
	if err := c.Bind(&product); err != nil {
		return err
	}
	if _, err := h.repo.CreateProduct(product); err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK,
		product.ID,
	)
}

func (h *Handler) GetAllProducts(c echo.Context) error {
	products, err := h.repo.GetAll()
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
