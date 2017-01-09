package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tutley/chiact/helpers"
	"github.com/tutley/chiact/models"
)

type userdata struct {
	FirstName string `json:"first"`
	LastName  string `json:"last"`
	Email     string `json:"email"`
	Password  string `json:"pass"`
}

// SignUpHandler is a Handler function for handling a user
// user signup route
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	db := startup.GetDb(r.Context())
	if db == nil {
		http.Error(w, "No database context", 500)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	userInfo := userdata{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	_, err = users.NewUser(userInfo.Email, userInfo.Password,
		userInfo.FirstName, userInfo.LastName, db)
	if err != nil {
		http.Error(w, err.Error(), 409)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"result":"success"}`))
}

// SignInHandler will return a JWT token for the user that signed in.
// This route must use the BasicMiddleware for authentication
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	db := startup.GetDb(r.Context())
	if db == nil {
		http.Error(w, "No database context", 500)
		return
	}

	user := startup.GetUser(r.Context())
	if user == nil {
		http.Error(w, "No user context", 401)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * 3600 * 24).Unix(),
	})

	tokenString, err := token.SignedString(startup.JwtSecret)
	if err != nil {
		http.Error(w, "Failed to create token", 401)
		return
	}

	fmt.Fprintf(w, "{\"token\": \"%s\"}", tokenString)
}
