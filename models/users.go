package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	User struct {
		Id	bson.ObjectId	`json:"id" bson:"_id,omitempty"`
		Username	string	`json:"username,omitempty" bson:"username,omitempty"`
		Firstname	string	`json:"firstname,omitempty" bson:"firstname,omitempty"`
		Lastname	string `json:"lastname,omitempty" bson:"lastname,omitempty"`
		Password	string	`json:"password,omitempty" bson:"password,omitempty"`
	}
)

func (u *User)IsNotDuplicate(mongoSession *mgo.Session) bool {
	err := mongoSession.DB("mdb").C("users").Find(bson.M{"username": u.Username}).One(u)
	if err != nil {
		return true
	}
	return false
}

func (u *User)SaveUserToDB(mongoSession *mgo.Session) error {
	err := mongoSession.DB("mdb").C("users").Insert(&u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User)ReadUsersFromDB(mongoSession *mgo.Session) ([]User , error ){
	result := []User{}
	err := mongoSession.DB("mdb").C("users").Find(nil).All(&result)
	if err != nil {
		return nil,err
	}
	return result, nil
}

func (u *User)ReadUsersByIDFromDB(mongoSession *mgo.Session) (*User, error) {
	err := mongoSession.DB("mdb").C("users").Find(bson.M{"_id": u.Id}).One(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User)UpdateUserToDB(mongoSession *mgo.Session) bool{
	err := mongoSession.DB("mdb").C("users").UpdateId(u.Id,u)
	if err != nil {
		return false
	}
	return true
}

func (u *User)DeleteUserByIDFromDB(mongoSession *mgo.Session) bool {
	err := mongoSession.DB("mdb").C("users").RemoveId(u.Id)
	if err != nil {
		return false
	}
	return true
}
