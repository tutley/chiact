package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/bradialabs/shortid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// User structure for working with users
type User struct {
	ID       string      `json:"id" bson:"_id,omitempty"`
	Email    string      `json:"email"`
	Password string      `json:"-"`
	Created  time.Time   `json:"created"`
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
}

// NewUser creates a new User and saves it in the database
func NewUser(email string, pass string,
	name string, db *mgo.Database) (*User, error) {
	//verify the email isn't already being used
	user, _ := FindUserByEmail(email, db)
	if user != nil {
		return nil, errors.New("stack|user: The user already exists")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("stack|user: %+v", err)
	}

	sid := shortid.New()
	sid.SetSeed(200)

	newUser := User{
		ID:       sid.Generate(),
		Email:    email,
		Password: string(passHash),
		Created:  time.Now(),
		Name:     name,
		Data:     nil,
	}

	err = newUser.Save(db)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// FindUserByID searches for an existing user with the passed ID
func FindUserByID(id string, db *mgo.Database) (*User, error) {
	user := User{}
	err := db.C("users").Find(bson.M{"_id": id}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByEmail searches for an existing user with the passed Email
func FindUserByEmail(email string, db *mgo.Database) (*User, error) {
	user := User{}
	err := db.C("users").Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Save Upserts the user into the database
func (user *User) Save(db *mgo.Database) error {
	_, err := db.C("users").UpsertId(user.ID, user)
	return err
}

// CheckPassword will check a passed password string with the stored hash
func (user *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
