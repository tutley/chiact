package helpers

import (
	"net/http"
	"gopkg.in/mgo.v2"
	"context"
)

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


