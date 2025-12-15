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

	r.Route("/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/{id}", handler.GetUserByID)
		r.Post("/email", handler.GetUserByEmail)
	})
}
