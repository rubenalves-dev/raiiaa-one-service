package products

import (
	"log"
	"net/http"

	"github.com/rubenalves-dev/raiiaa-one-service/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 2. Return JSON in an HTTP Request

	products := []string{"Product 1", "Product 2", "Product 3"}

	json.Write(w, http.StatusOK, products)
}
