package repository

import (
	"bootcamp-auth-microservice/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	GetUsersByID(id string) (*models.User, error)
	ReadUser(filter models.UserFilter, page, size int) ([]models.UserView, error)
	StudentRegister(user *models.User) error
	UpdateName(id string, user *models.UpdateName) (*models.UpdateName, error)
}

func (r *RepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT * FROM bootcamp_users WHERE username = ?"

	var user models.User
	err := r.DB.Read.Get(&user, query, username)
	if err != nil {
		log.Error().Err(err).Msg("Something went wrong")
		return nil, err
	}
	return &user, nil
}

func (r *RepositoryImpl) StudentRegister(user *models.User) error {
	query :=
		`
	INSERT INTO bootcamp_users (id, username, name, role, password, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	user.ID = uuid.New().String()
	// Hash the user's password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Role = "student"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err = r.DB.Write.Exec(
		query,
		user.ID,
		user.Username,
		user.Name,
		user.Role,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryImpl) UpdateName(id string, user *models.UpdateName) (*models.UpdateName, error) {
	query :=
		`
	UPDATE bootcamp_users 
	SET name = ?, updated_at = ? 
	WHERE id = ?
	`
	user.ID = id
	user.UpdatedAt = time.Now()
	_, err := r.DB.Write.Exec(
		query,
		user.Name,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *RepositoryImpl) ReadUser(filter models.UserFilter, page, size int) ([]models.UserView, error) {
	query := "SELECT username, name, role, created_at FROM bootcamp_users"

	// Add filters
	args := []interface{}{}
	if filter.Name != "" {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+filter.Name+"%")
	}

	// Add pagination
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	query += " LIMIT ? OFFSET ?"
	offset := (page - 1) * size
	args = append(args, size, offset)

	var users []models.UserView
	err := r.DB.Read.Select(&users, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Something went wrong")
		return nil, err
	}
	return users, nil
}

func (r *RepositoryImpl) GetUsersByID(id string) (*models.User, error) {
	query := "SELECT * FROM bootcamp_users WHERE id = ?"

	var variant models.User
	err := r.DB.Read.Get(&variant, query, id)
	if err != nil {
		log.Error().Err(err).Msg("Something went wrong")
		return nil, err
	}
	return &variant, nil
}
