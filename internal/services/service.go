package services

import (
	"bootcamp-auth-microservice/internal/models"
	"bootcamp-auth-microservice/internal/repository"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	StudentRegister(user *models.User) error
	GenerateJWT(user *models.User) (string, error)
	Login(username, password string) (string, error)
	ReadUser(filter models.UserFilter, page, size int) ([]models.UserView, error)
	UpdateName(id string, user *models.UpdateName) (*models.UpdateName, error)
	GetUsersByID(id string) (*models.User, error)
	ValidateJWT(tokenString string) (*models.User, error)
}

type ServiceImpl struct {
	Repo repository.Repository
}

func ProvideService(r repository.Repository) *ServiceImpl {
	return &ServiceImpl{
		Repo: r,
	}
}

func (s *ServiceImpl) StudentRegister(user *models.User) error {
	return s.Repo.StudentRegister(user)
}

func (s *ServiceImpl) UpdateName(id string, user *models.UpdateName) (*models.UpdateName, error) {
	return s.Repo.UpdateName(id, user)
}

func (s *ServiceImpl) ReadUser(filter models.UserFilter, page, size int) ([]models.UserView, error) {
	return s.Repo.ReadUser(filter, page, size)
}

func (s *ServiceImpl) GetUsersByID(id string) (*models.User, error) {
	return s.Repo.GetUsersByID(id)
}

func (s *ServiceImpl) UserCheck(username, password string) (*models.User, error) {
	// Fetch the user by username from the repository
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	// Check if the user exists and the password is correct
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, models.ErrUnauthorized
	}

	return user, nil
}

func (s *ServiceImpl) GenerateJWT(user *models.User) (string, error) {
	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// In a production environment, keep the secret key secure and don't hardcode it here.
	secretKey := []byte("your-secret-key")

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *ServiceImpl) Login(username, password string) (string, error) {

	user, err := s.UserCheck(username, password)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *ServiceImpl) ValidateJWT(tokenString string) (*models.User, error) {
	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, jwt.ErrSignatureInvalid
		}
		// In a production environment, keep the secret key secure and don't hardcode it here.
		secretKey := []byte("your-secret-key")
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract the user information from the claims
		userID := claims["user_id"].(string)
		username := claims["username"].(string)
		role := claims["role"].(string)

		// Create a user model with the extracted information
		user := &models.User{
			ID:       userID,
			Username: username,
			Role:     role,
		}

		return user, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
