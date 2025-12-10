package users

import (
	"log"
	"net/http"

	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
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

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body repo.CreateUserParams
	if err := json.Read(r, &body); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusCreated, user)
}
