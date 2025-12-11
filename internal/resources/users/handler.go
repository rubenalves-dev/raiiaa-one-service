package users

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/common/validators"
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

func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	uuid, err := uuid.Parse(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(r.Context(), uuid)
	if err != nil {
		statusCode := http.StatusInternalServerError
		log.Println(err)
		if errors.Is(err, ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	json.Write(w, http.StatusOK, user)
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
		statusCode := http.StatusInternalServerError
		log.Println(err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	json.Write(w, http.StatusCreated, user)
}

func (h *handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	var body GetUserByEmailParams
	if err := json.Read(r, &body); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		statusCode := http.StatusInternalServerError
		log.Println(err)
		if errors.Is(err, validators.ErrInvalidEmail) {
			statusCode = http.StatusBadRequest
		}
		if errors.Is(err, ErrUserNotFound) {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	json.Write(w, http.StatusOK, user)
}
