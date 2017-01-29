package models

import (
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Page struct {
	ID string `json:"id" bson:"_id,omitempty"`
	URL string `json:"url"`
	Modified time.Time `json:"modified"`
	Content string `json:"content"`
}

// FindUserByEmail searches for an existing user with the passed Email
func FindPageByURL(url string, db *mgo.Database) (*Page, error) {
	page := Page{}
	err := db.C("pages").Find(bson.M{"url": url}).One(&page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Save Upserts the user into the database
func (page *Page) Save(db *mgo.Database) error {
	_, err := db.C("pages").Upsert(bson.M{"url": page.URL}, page)
	return err
}