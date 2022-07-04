package repository

import "813-online-cinema-monolit/pkg/models"

type Db interface {
	GetUserById(id uint64) (*models.User, error)
	GetUserByLoginPassword(login, password string) (*models.User, error)
	GetMovie(id uint64) (*models.Movie, error)
	GetUsersMovie(id uint64) ([]*models.Movie, error)
}
