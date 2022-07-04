package server

import (
	"813-online-cinema-monolit/pkg/models"
	"813-online-cinema-monolit/pkg/repository/sqlite"
	"813-online-cinema-monolit/pkg/templateInitialization"
	"context"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	templatesFolder        = "templates"
	get             string = "GET"
	post            string = "POST"
	templates, _           = templateInitialization.LoadTemplates(templatesFolder)
)

type Route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func NewRoute(method, pattern string, handler http.HandlerFunc) Route {
	return Route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func (oc *OnlineCinema) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range oc.routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func UserHandler(t *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write([]byte("Welcome, user " + getField(request, 0))); err != nil {
			log.Println(err.Error())
		}
	}
}

func GetHandler(t *template.Template, data any) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := t.Execute(writer, data); err != nil {
			log.Fatal(err)
		}
	}
}

func GetHandlerMovie(t *template.Template) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.ParseInt(getField(request, 0), 10, 64)
		if err != nil {
			log.Println(err)
			return
		}
		if err = t.Execute(writer, struct {
			Id int64
		}{Id: id}); err != nil {
			log.Fatal(err)
		}
	}
}

func LoginPostHandler(db *sqlite.Sqlite) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			user   *models.User
			movies []*models.Movie
			err    error
		)
		if err = request.ParseForm(); err != nil {
			log.Println(err)
		}
		user, err = db.GetUserByLoginPassword(request.Form["login"][0], request.Form["password"][0])
		if err != nil {
			if err = templates["login"].Execute(writer, nil); err != nil {
				log.Println(err)
				return
			}
		}
		movies, err = db.GetUsersMovie(user.Id)
		if err != nil {
			log.Println(err)
			return
		}
		AuthorizedUser.User = user
		AuthorizedUser.Movies = movies

		http.Redirect(writer, request, "/list", http.StatusFound)
		/*if err = templates["list"].Execute(writer, struct {
			User   *models.User
			Movies []*models.Movie
		}{User: user,
			Movies: movies}); err != nil {
			log.Println(err)
			return
		}*/
	}
}
