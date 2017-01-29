package helpers

import (
	"log"
	"net/http"
	"context"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dgrijalva/jwt-go"
	"github.com/tutley/chiact/models"
)

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
		user, err := models.FindUserByID(claims["id"].(string), db)
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


