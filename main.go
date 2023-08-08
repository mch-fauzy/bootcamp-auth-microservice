package main

import (
	"bootcamp-auth-microservice/infras"
	"bootcamp-auth-microservice/internal/repository"
	"bootcamp-auth-microservice/internal/services"
	"bootcamp-auth-microservice/transport/middleware"
	"bootcamp-auth-microservice/transport/routes"
	"fmt"
	"net/http"
)

func main() {
	secretKey := []byte("your-secret-key")
	// Create a new database connection
	db := infras.ProvideConn()

	// Initialize the repository with the database connection
	repo := repository.ProvideRepo(&db)

	// Initialize the service with the repository
	svc := services.ProvideService(repo)

	// Initialize the authentication middleware

	auth := middleware.ProvideAuthentication(&db, secretKey)

	// Initialize the router with the service and authentication
	r := routes.ProvideRouter(svc, auth)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r.SetupRoutes())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
