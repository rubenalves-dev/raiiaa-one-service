package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
)

func InitRoutes(r *chi.Mux, db *pgx.Conn) {
	repo := repo.New(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.Post("/users", handler.CreateUser)
	r.Get("/users/{id}", handler.GetUserByID)
	r.Post("/users/email", handler.GetUserByEmail)
}
