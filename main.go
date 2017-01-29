package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/tutley/chiact/handlers"
	"github.com/tutley/chiact/helpers"
	"time"
	"io/ioutil"
)

var index []byte

func main() {
	helpers.SetJwtSecret([]byte("youshouldchangethissecret"))

	// setup the Chi router
	r := chi.NewRouter()

	// Use Chi built-in middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Mount database
	r.Use(helpers.MongoMiddleware)

	// When a client closes their connection midway through a request, the
	// http.CloseNotifier will cancel the request context (ctx).
	r.Use(middleware.CloseNotify)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Setup routes/routers for the API. The routers are defined last in this file
	r.Post("/api/1/signup", handlers.SignUpHandler)
	r.Mount("/api/1", APIRouter())
	r.Mount("/api/1/login", LoginRouter())

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "client")

	// load the index html
	index, _ = ioutil.ReadFile(filesDir+"/index.html")

	// got-dang favicon
	r.Get("/favicon.ico", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filesDir+"/favicon.ico")
	})

	// Expose the client directory for assets
	r.FileServer("/client", http.Dir(filesDir))

	// Serve the index
	r.Mount("/", RootRouter())

	// and.... go!
	http.ListenAndServe(":3333", r) // TODO: Make server port a config var
	// TODO: ALso make server port variable the same thing used in the prerender middleware
}

// ROUTING
// This section defines the routes that will be served

// LoginRouter provides the routes for loging in
func LoginRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(helpers.BasicMiddleware)

	r.Get("/", handlers.SignInHandler)
	return r
}

// APIRouter handles the routes for the API
func APIRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(helpers.JwtAuthMiddleware)

	r.Get("/me", handlers.GetMeHandler)
	r.Put("/me", handlers.UpdateMeHandler)

	return r
}

func RootRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(helpers.PrerenderMiddleware)

	// catch any remaining routes and serve them the index.html
	// let react-router deal with them
	// TODO: handle 404 in react router http://knowbody.github.io/react-router-docs/guides/NotFound.html
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write(index)
	})

	return r
}