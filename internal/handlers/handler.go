package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"bootcamp-auth-microservice/internal/models"
	"bootcamp-auth-microservice/internal/services"

	"github.com/go-chi/chi"
)

type Handler struct {
	Service services.Service
}

func ProvideHandler(service services.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) StudentRegister(w http.ResponseWriter, r *http.Request) {
	// Define the required struct for the request body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password fields are required", http.StatusBadRequest)
		return
	}

	//Might need to check if username and uuid is exist or not (to avoid duplicate)

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}

	err := h.Service.StudentRegister(user)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "User successfully registered",
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateName(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	// Check if the user with the given ID exists in the database
	_, err := h.Service.GetUsersByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	// Define the required struct for the request body
	var req struct {
		Name string `json:"name"`
	}

	// Decode the request body into the req struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" {
		http.Error(w, "name field is required", http.StatusBadRequest)
		return
	}

	// next may add validation "teacher" cant update other teacher profile

	name := &models.UpdateName{
		Name: req.Name,
	}

	updatedName, err := h.Service.UpdateName(userID, name)
	if err != nil {
		http.Error(w, "Failed to update variant", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "User successfully updated",
		"user":    updatedName,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ReadUser(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	q := r.URL.Query()
	name := q.Get("name")
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("size"))

	// Set default values for page and size if they are not provided or invalid
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	// Call the service to fetch the user with the specified filters and sorting
	users, err := h.Service.ReadUser(models.UserFilter{
		Name: name,
	},
		page,
		size)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) ValidateAuth(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Validate the JWT token
	user, err := h.Service.ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Respond with the user information
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the username and password from the request body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password fields are required", http.StatusBadRequest)
		return
	}

	// Authenticate user and generate JWT token
	token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		if err == models.ErrUnauthorized {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to authenticate user", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the JWT token
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	response := map[string]string{
		"token": token,
	}
	json.NewEncoder(w).Encode(response)
}
