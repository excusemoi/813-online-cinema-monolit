package sqlite

import (
	"813-online-cinema-monolit/pkg/models"
	"database/sql"
	"errors"
	"fmt"
)

var (
	createUserTableQuery = `create table if not exists user (id integer not null primary key,
login text, password text, name text, surname text);`
	createMovieTableQuery = `create table if not exists movie (id integer not null primary key,
name text, description text);`
	createUserMovieTableQuery      = `create table if not exists user_movie (id_user, id_movie);`
	selectUserByIdQuery            = `select * from user where id = ?;`
	selectUserByLoginPasswordQuery = `select * from user where login = ? and password = ?;`
	selectMovieByIdQuery           = `select * from movie where id = ?;`
	selectUserMoviesByIdQuery      = `select id_movie from user_movie where id_user = ?;`
	insertUserQuery                = `insert into user values(?,?,?,?,?);`
	insertMovieQuery               = `insert into movie values(?,?,?);`
	insertUserMovieQuery           = `insert into user_movie values(?,?);`
)

type Sqlite struct {
	db *sql.DB
}

func (s *Sqlite) InitDb() error {
	var err error
	if s.db == nil {
		return errors.New("db uninitialized")
	}
	if _, err = s.db.Exec(createUserTableQuery); err != nil {
		return err
	}
	if _, err = s.db.Exec(createMovieTableQuery); err != nil {
		return err
	}
	if _, err = s.db.Exec(createUserMovieTableQuery); err != nil {
		return err
	}
	return nil
}

func NewSqlite(db *sql.DB) *Sqlite {
	return &Sqlite{db: db}
}

func (s *Sqlite) InsertUser(user *models.User) error {
	if _, err := s.db.Exec(insertUserQuery, user.Id, user.Login, user.Password, user.Name, user.Surname); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) InsertMovie(movie *models.Movie) error {
	if _, err := s.db.Exec(insertMovieQuery, movie.Id, movie.Name, movie.Description); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) InsertUserMovie(userMovie *models.UserMovie) error {
	if _, err := s.db.Exec(insertUserMovieQuery, userMovie.IdUser, userMovie.IdMovie); err != nil {
		return err
	}
	return nil
}

func (s *Sqlite) GetUserById(id uint64) (*models.User, error) {
	row := s.db.QueryRow(selectUserByIdQuery, id)
	user := models.User{}
	if err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Name, &user.Surname); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Sqlite) GetUserByLoginPassword(login, password string) (*models.User, error) {
	row := s.db.QueryRow(selectUserByLoginPasswordQuery, login, password)
	if row == nil {
		return nil, errors.New(fmt.Sprintf("user with login %s unregistered", login))
	}
	user := models.User{}
	if err := row.Scan(&user.Id, &user.Login, &user.Password, &user.Name, &user.Surname); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Sqlite) GetMovie(id uint64) (*models.Movie, error) {
	row := s.db.QueryRow(selectMovieByIdQuery, id)
	movie := models.Movie{}
	if err := row.Scan(&movie.Id, &movie.Name, &movie.Description); err != nil {
		return nil, err
	}
	return &movie, nil
}
func (s *Sqlite) GetUsersMovie(id uint64) ([]*models.Movie, error) {
	rows, err := s.db.Query(selectUserMoviesByIdQuery, id)
	movies := make([]*models.Movie, 0)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			movieId uint64
			movie   *models.Movie
		)
		if err = rows.Scan(&movieId); err != nil {
			return nil, err
		}
		if movie, err = s.GetMovie(movieId); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}
