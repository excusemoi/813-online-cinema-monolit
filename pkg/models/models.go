package models

type Movie struct {
	Id          uint64
	Name        string
	Description string
}

type User struct {
	Id       uint64
	Login    string
	Password string
	Name     string
	Surname  string
}

type UserMovie struct {
	IdUser  uint64
	IdMovie uint64
}

type AuthorizedUser struct {
	User   *User
	Movies []*Movie
}
