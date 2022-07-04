package server

import (
	"813-online-cinema-monolit/pkg/models"
	"813-online-cinema-monolit/pkg/repository"
	"log"
	"net/http"
)

var (
	AuthorizedUser = &models.AuthorizedUser{}
)

type OnlineCinema struct {
	db     repository.Db
	server *http.Server
	routes []Route
}

func NewOnlineCinema(db repository.Db, routes []Route) *OnlineCinema {
	return &OnlineCinema{db: db, routes: routes}
}

func (oc *OnlineCinema) Run(templateFolder string) error {
	var (
		err error
	)

	oc.server = &http.Server{
		Addr:    ":8000",
		Handler: oc,
	}

	log.Println("Server listening on", oc.server.Addr)
	if err = oc.server.ListenAndServe(); err != nil {
		return err
	}

	return err
}
