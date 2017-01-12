package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"context"
	"encoding/base64"
	"strings"

	mgo "gopkg.in/mgo.v2"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/tutley/chiact/models"

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

	// Mount database
	r.Use(ChiMongoMiddleware)

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

// This section contains helper functions and local middlewares
// TODO : figure out a way to separate this out into a helper function file

// ChiMongoMiddleware gives us our connection to the database
func ChiMongoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// setup the mgo connection
		session, err := mgo.Dial("mongodb://localhost/") // TODO: make this a config var

		if err != nil {
			panic(err)
		}

		reqSession := session.Clone()
		defer reqSession.Close()
		db := reqSession.DB("chiact") // TODO: Make this a config var
		ctx := context.WithValue(r.Context(), startup.DbKey, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ChiBasicMiddleware lets us process logins
func ChiBasicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := startup.GetDb(r.Context())
		if db == nil {
			log.Print("No database context")
			http.Error(w, "Not authorized", 401)
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		//Find the user in the database
		user, err := users.FindUserByEmail(pair[0], db)
		if err != nil || user == nil {
			log.Printf("User %+v not found.", pair[0])
			http.Error(w, "Not authorized", 401)
			return
		}

		//Check their password
		err = user.CheckPassword(pair[1])
		if err != nil {
			log.Printf("Invalid password for User: %+v.", pair[0])
			http.Error(w, "Not authorized", 401)
			return
		}

		//clear password
		user.Password = ""

		//Set the logged in user in the context
		ctx := context.WithValue(r.Context(), startup.UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// ChiJwtAuthMiddleware handles the JWT authentication strategy
func ChiJwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := startup.GetDb(r.Context())
		if db == nil {
			log.Print("No database context")
			http.Error(w, "Not authorized", 401)
		}

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return startup.JwtSecret, nil
		})

		if err != nil {
			http.Error(w, "Invalid Token", 401)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid Token", 401)
			return
		}

		// Check signing method to avoid vulnerabilities
		if token.Method != jwt.SigningMethodHS256 {
			http.Error(w, "Invalid Token", 401)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		//Find the user in the database
		user, err := users.FindUserByID(claims["id"].(string), db)
		if err != nil || user == nil {
			log.Printf("User %s not found.", claims["id"].(string))
			http.Error(w, "Not authorized", 401)
			return
		}

		//clear password
		user.Password = ""

		//Set the logged in user in the context
		ctx := context.WithValue(r.Context(), startup.UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ROUTING
// This section defines the routes that will be served

// LoginRouter provides the routes for loging in
func LoginRouter() chi.Router {
	r := chi.NewRouter()

	//r.Use(ChiMongoMiddleware)
	r.Use(ChiBasicMiddleware)

	r.Get("/", auth.SignInHandler)
	return r
}

// APIRouter handles the routes for the API
func APIRouter() chi.Router {
	r := chi.NewRouter()

	//r.Use(ChiMongoMiddleware)
	r.Use(ChiJwtAuthMiddleware)

	r.Get("/me", auth.GetMeHandler)
	r.Put("/me", auth.UpdateMeHandler)

	return r
}
