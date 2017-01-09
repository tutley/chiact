package helpers

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	mgo "gopkg.in/mgo.v2"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/tutley/chiact/models"
)

type key int

var dbKey key = 100000
var userKey key = 200000

var JwtSecret []byte

// GetDb grabs the mgo database from the context
func GetDb(ctx context.Context) *mgo.Database {
	return ctx.Value(dbKey).(*mgo.Database)
}

// GetUser grabs the current user from the context
func GetUser(ctx context.Context) *users.User {
	return ctx.Value(userKey).(*users.User)
}

// SetJwtSecret sets the secret that will be used to sign and verify JWT tokens
func SetJwtSecret(secret []byte) {
	JwtSecret = secret
}

// MongoMiddleware adds mgo MongoDB to context
func MongoMiddleware(dbName string, connectURI string) func(next http.Handler) http.Handler {
	// setup the mgo connection
	session, err := mgo.Dial(connectURI)

	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqSession := session.Clone()
			defer reqSession.Close()
			db := reqSession.DB(dbName)
			ctx := context.WithValue(r.Context(), dbKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// BasicMiddleware adds Basic Auth to routes. This assumes
// that the route is using the MongoMiddleware
func BasicMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
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
		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

// JwtAuthMiddleware JWT Bearer Authentication to routes. This assumes
// that the route is using the MongoMiddleware
func JwtAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
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
		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
