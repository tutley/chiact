package startup

import (
	"context"

	mgo "gopkg.in/mgo.v2"

	"github.com/tutley/chiact/models"
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
