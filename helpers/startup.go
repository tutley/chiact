package helpers

import (
	"context"

	mgo "gopkg.in/mgo.v2"

	"github.com/tutley/chiact/models"
	"net/http"
	"log"
	"strings"

	"encoding/base64"

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
func GetUser(ctx context.Context) *models.User {
	return ctx.Value(UserKey).(*models.User)
}

// SetJwtSecret sets the secret that will be used to sign and verify JWT tokens
func SetJwtSecret(secret []byte) {
	JwtSecret = secret
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
		user, err := models.FindUserByEmail(pair[0], db)
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

