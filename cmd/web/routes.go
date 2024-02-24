package main

import (
	"net/http"
	"path/filepath"


	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	// mux := http.NewServeMux()
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	// mux.Handle("/static", http.NotFoundHandler())
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("/", app.home)
	router.HandlerFunc(http.MethodGet, "/", app.home)
	// mux.HandleFunc("/snippet/view", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	// mux.HandleFunc("/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateForm)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreate)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}	

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()

			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}

	}

	return f, nil
}
