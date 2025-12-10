package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/ignore/products"
)

type application struct {
	config config
	// db driver
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string // data source name
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)                 // rate limiting
	r.Use(middleware.RealIP)                    // rate limiting, analytics and tracing
	r.Use(middleware.Logger)                    // better logs
	r.Use(middleware.Recoverer)                 // recover from crashes
	r.Use(middleware.Timeout(60 * time.Second)) // stops process on request timeout

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("The server is healthy"))
	})

	productsService := products.NewService()
	productsHandlers := products.NewHandler(productsService)
	r.Get("/products", productsHandlers.ListProducts)

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
