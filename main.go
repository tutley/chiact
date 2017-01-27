package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/tutley/chiact/handlers"
	"github.com/tutley/chiact/helpers"
)

func main() {
	startup.SetJwtSecret([]byte("youshouldchangethissecret"))

	r := chi.NewRouter()

	// Use Chi built-in middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Mount database
	r.Use(startup.MongoMiddleware)

	// Serve the client
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "client")
	r.FileServer("/", http.Dir(filesDir))

	// Setup routes/routers for the API. The routers are defined last in this file
	r.Post("/api/1/signup", auth.SignUpHandler)
	r.Mount("/api/1", APIRouter())
	r.Mount("/api/1/login", LoginRouter())

	// and.... go!
	http.ListenAndServe(":3333", r) // TODO: Make this port a config var
}

// ROUTING
// This section defines the routes that will be served

// LoginRouter provides the routes for loging in
func LoginRouter() chi.Router {
	r := chi.NewRouter()

	//r.Use(ChiMongoMiddleware)
	r.Use(startup.BasicMiddleware)

	r.Get("/", auth.SignInHandler)
	return r
}

// APIRouter handles the routes for the API
func APIRouter() chi.Router {
	r := chi.NewRouter()

	//r.Use(ChiMongoMiddleware)
	r.Use(startup.JwtAuthMiddleware)

	r.Get("/me", auth.GetMeHandler)
	r.Put("/me", auth.UpdateMeHandler)

	return r
}
