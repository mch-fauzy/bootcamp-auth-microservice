package routes

import (
	"bootcamp-auth-microservice/internal/services"
	custom_middleware "bootcamp-auth-microservice/transport/middleware"
	"net/http"

	"bootcamp-auth-microservice/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	Handler        *handlers.Handler
	Authentication *custom_middleware.Authentication
}

func ProvideRouter(service services.Service, auth *custom_middleware.Authentication) *Router {
	handler := handlers.ProvideHandler(service)
	return &Router{
		Handler:        handler,
		Authentication: auth,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Route("/v1", func(rc chi.Router) {
		// Use the authentication middleware for the protected endpoints
		rc.Use(r.Authentication.VerifyJWT)

		rc.Get("/validate-auth", r.Handler.ValidateAuth)

		// Protected endpoints accessible only by teachers
		rc.Group(func(rc chi.Router) {
			rc.Use(r.Authentication.VerifyTeacherJWT)
			rc.Get("/users", r.Handler.ReadUser)
			rc.Put("/users/{id}", r.Handler.UpdateName)
		})
	})

	// Public endpoints
	mux.Post("/v1/login", r.Handler.Login)
	mux.Post("/v1/register", r.Handler.StudentRegister)
	return mux
}
