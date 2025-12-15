package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
	refreshtokens "github.com/rubenalves-dev/raiiaa-one-service/internal/resources/refresh-tokens"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/resources/users"
)

func InitRoutes(r *chi.Mux, db *pgx.Conn) {
	repo := repo.New(db)
	usersService := users.NewService(repo)
	tokensService := refreshtokens.NewService(repo)
	service := NewService(usersService, tokensService)
	handler := NewHandler(service)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}
