package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/tutley/chiact/handlers"
	"github.com/tutley/chiact/helpers"
	"golang.org/x/net/context"
)

type UserData struct {
	Stuff string
}

func main() {
	helpers.SetJwtSecret([]byte("secret"))

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/api/1/signup", auth.SignUpHandler)
	r.Mount("/api/1", ApiRouter())
	r.Mount("/api/1/login", LoginRouter())

	http.ListenAndServe(":3333", r)
}

func ChiMongoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(helpers.MongoMiddleware("chitest", "", next))
}

func ChiBasicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(helpers.BasicMiddleware(next))
}

func ChiJwtAuthMiddleware(next chi.Handler) http.Handler {
	return http.HandlerFunc(helpers.JwtAuthMiddleware(next))
}

func LoginRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(ChiMongoMiddleware)
	r.Use(ChiBasicMiddleware)

	r.Get("/", auth.SignInHandler)
	return r
}

func ApiRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(ChiMongoMiddleware)
	r.Use(ChiJwtAuthMiddleware)

	r.Get("/me", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		user := helpers.GetUser(ctx)
		j, er := json.Marshal(&user)
		if er != nil {
			log.Fatal(er)
		}
		w.Write(j)
	})

	r.Put("/me", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		db := helpers.GetDb(ctx)
		user := helpers.GetUser(ctx)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		userData := UserData{}
		err = json.Unmarshal(body, &userData)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		user.Data = userData
		user.Save(db)

		j, er := json.Marshal(&user)
		if er != nil {
			log.Fatal(er)
		}
		w.Write(j)
	})

	return r
}
