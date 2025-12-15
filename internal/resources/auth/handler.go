package auth

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

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var req repo.CreateUserParams

	if err := json.Read(r, &req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payoad", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.Write(w, http.StatusCreated, user)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginParams
	if err := json.Read(r, &req); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	json.Write(w, http.StatusOK, map[string]string{"access_token": token})
}
