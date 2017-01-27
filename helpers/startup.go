package startup

import (
	"context"

	mgo "gopkg.in/mgo.v2"

	"github.com/tutley/chiact/models"
	"net/http"
	"log"
	"strings"

	"encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type key int

// DbKey is somethign I need to document
var DbKey key = 100000

// UserKey is something I need to document
var UserKey key = 200000

// JwtSecret is something I Need ot document
var JwtSecret []byte

// GetDb grabs the mgo database from the context
func GetDb(ctx context.Context) *mgo.Database {
	return ctx.Value(DbKey).(*mgo.Database)
}

// GetUser grabs the current user from the context
func GetUser(ctx context.Context) *users.User {
	return ctx.Value(UserKey).(*users.User)
}

// SetJwtSecret sets the secret that will be used to sign and verify JWT tokens
func SetJwtSecret(secret []byte) {
	JwtSecret = secret
}


// MongoMiddleware gives us our connection to the database
func MongoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// setup the mgo connection
		session, err := mgo.Dial("mongodb://localhost/") // TODO: make this a config var

		if err != nil {
			panic(err)
		}

		reqSession := session.Clone()
		defer reqSession.Close()
		db := reqSession.DB("chiact") // TODO: Make this a config var
		ctx := context.WithValue(r.Context(), DbKey, db)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


// BasicMiddleware lets us process logins
func BasicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := GetDb(r.Context())
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
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// JwtAuthMiddleware handles the JWT authentication strategy
func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db := GetDb(r.Context())
		if db == nil {
			log.Print("No database context")
			http.Error(w, "Not authorized", 401)
		}

		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
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
		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


