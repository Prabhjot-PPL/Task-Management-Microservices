// package routes

// import (
// 	"net/http"
// 	"task_service/src/internal/interfaces/input/api/rest/handler"

// 	"github.com/go-chi/chi/v5"
// )

// func InitRoutes(handler handler.TaskHandler) http.Handler {

// 	r := chi.NewRouter()

// 	r.Post("/tasks", handler.CreateTaskHandler)
// 	r.Put("/tasks/{id}", handler.UpdateTaskHandler)
// 	r.Get("/tasks", handler.ListTasksHandler)

// 	return r
// }

package routes

import (
	"fmt"
	"net/http"
	sessionclient "task_service/src/internal/adaptors/grpcclient"
	"task_service/src/internal/interfaces/input/api/rest/handler"

	"github.com/go-chi/chi/v5"
)

func SessionAuthMiddleware(sessionClient *sessionclient.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id") // <-- change this to your cookie name
			if err != nil {
				http.Error(w, "session cookie is missing", http.StatusUnauthorized)
				return
			}
			sessionID := cookie.Value
			fmt.Println("session id : ", sessionID)

			valid, _, err := sessionClient.ValidateSession(r.Context(), sessionID)
			if err != nil || !valid {
				http.Error(w, "invalid session", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func InitRoutes(taskHandler handler.TaskHandler, sessionClient *sessionclient.Client) http.Handler {
	r := chi.NewRouter()

	// Apply session middleware on task routes
	r.Route("/tasks", func(r chi.Router) {
		r.Use(SessionAuthMiddleware(sessionClient)) // protect all /tasks routes

		r.Post("/", taskHandler.CreateTaskHandler)
		r.Put("/{id}", taskHandler.UpdateTaskHandler)
		// r.Get("/{id}", taskHandler.GetTask)
		// ... other task routes
	})

	return r
}
