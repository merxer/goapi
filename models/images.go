package models

import (
	_ "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	Image struct {
		Id	bson.ObjectId	`json:"id,omitempty" bson:"_id,omitempty"`
		Base64String	string	`json:"base64string,omitempty" bson:"base64string,omitempty"`
		Name string	`json:"name,omitempty" bson:"name,omitempty"`
		Url string	`json:"url,omitempty" bson:"url,omitempty"`
	}
)
