package routes

import (
	"net/http"
	"user_service/src/internal/interfaces/input/api/rest/handler"
	pkgmiddleware "user_service/src/internal/interfaces/input/api/rest/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitRoutes(handler handler.UserHandler) http.Handler {

	r := chi.NewRouter()

	// Build-in middlware methods from chi
	r.Use(middleware.Logger)

	r.Post("/auth/register", handler.RegisterHandler)

	r.Post("/auth/login", handler.LoginHandler)

	// Protected routes
	r.Group(func(protected chi.Router) {
		protected.Use(pkgmiddleware.SessionAuth)
		protected.Get("/user", handler.ProfileHandler)
		protected.Post("/user", handler.UpdateHandler)
	})

	return r
}
