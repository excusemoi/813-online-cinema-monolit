package main

import (
	"813-online-cinema-monolit/pkg/repository/sqlite"
	"813-online-cinema-monolit/pkg/server"
	"813-online-cinema-monolit/pkg/templateInitialization"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	dbFileName      = "cinema.db"
	dbDriverName    = "sqlite3"
	templatesFolder = "templates"
	templates, _    = templateInitialization.LoadTemplates(templatesFolder)
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	db, err := sql.Open(dbDriverName, dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	sqliteDb := sqlite.NewSqlite(db)
	if err = sqliteDb.InitDb(); err != nil {
		log.Fatal(err)
	}

	routes := []server.Route{
		server.NewRoute("GET", "/", server.GetHandler(templates["welcome"], nil)),
		server.NewRoute("GET", "/login", server.GetHandler(templates["login"], nil)),
		server.NewRoute("POST", "/login", server.LoginPostHandler(sqliteDb)),
		server.NewRoute("GET", "/list", server.GetHandler(templates["list"], server.AuthorizedUser)),
		server.NewRoute("GET", "/movie/([0-9]+)", server.GetHandlerMovie(templates["movie"])),
	}

	oc := server.NewOnlineCinema(sqliteDb, routes)

	go func() {
		if err = oc.Run(templatesFolder); err != nil {
			log.Fatal(err)
		}
	}()

	<-done

	log.Println("Server stopped")
}
