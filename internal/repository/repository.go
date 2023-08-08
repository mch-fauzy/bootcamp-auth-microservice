package repository

import "bootcamp-auth-microservice/infras"

type Repository interface {
	UserRepository
}

type RepositoryImpl struct {
	DB *infras.Conn
}

func ProvideRepo(db *infras.Conn) *RepositoryImpl {
	return &RepositoryImpl{
		DB: db,
	}
}
