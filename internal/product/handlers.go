package product

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo  *Repository
	kafka *KafkaProducer
}

func NewHandler(db *sql.DB, kafka *KafkaProducer) *Handler {
	return &Handler{
		repo:  NewRepository(db),
		kafka: kafka,
	}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreatePruduct(product); err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	h.kafka.SendProductMessage(&product)

	json.NewEncoder(w).Encode(product)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetById(id)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product := new(Product)
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	product.ID = id
	if err := h.repo.Update(product); err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	h.kafka.SendProductMessage(product)

	json.NewEncoder(w).Encode(product)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
