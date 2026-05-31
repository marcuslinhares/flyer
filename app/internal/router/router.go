package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/marcuslinhares/flyer/app/internal/auth"
	"github.com/marcuslinhares/flyer/app/internal/handlers"
	"github.com/marcuslinhares/flyer/app/internal/middleware"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.CORSMiddleware)
	r.Use(chimw.Recoverer)

	// Health check (no auth required)
	r.Get("/api/health", handlers.HealthCheck)

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(auth.APIKeyMiddleware)

		// Flyers CRUD
		r.Get("/flyers", handlers.ListFlyers)
		r.Post("/flyers", handlers.CreateFlyer)
		r.Get("/flyers/{id}", handlers.GetFlyer)
		r.Put("/flyers/{id}", handlers.UpdateFlyer)
		r.Delete("/flyers/{id}", handlers.DeleteFlyer)

		// Estabelecimentos CRUD
		r.Get("/estabelecimentos", handlers.ListEstabelecimentos)
		r.Post("/estabelecimentos", handlers.CreateEstabelecimento)
		r.Get("/estabelecimentos/{id}", handlers.GetEstabelecimento)
		r.Put("/estabelecimentos/{id}", handlers.UpdateEstabelecimento)
		r.Delete("/estabelecimentos/{id}", handlers.DeleteEstabelecimento)

		// Categorias CRUD
		r.Get("/categorias", handlers.ListCategorias)
		r.Post("/categorias", handlers.CreateCategoria)
		r.Get("/categorias/{id}", handlers.GetCategoria)
		r.Put("/categorias/{id}", handlers.UpdateCategoria)
		r.Delete("/categorias/{id}", handlers.DeleteCategoria)

		// Upload
		r.Post("/upload", handlers.UploadFile)
	})

	// Serve uploaded files (no auth)
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/app/uploads"
	}
	absDir, _ := filepath.Abs(uploadDir)
	fileServer := http.StripPrefix("/uploads/", http.FileServer(http.Dir(absDir)))
	r.Get("/uploads/*", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	return r
}
