package models

import (
	_ "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"../helper"
)

type (
	User struct {
		Id	bson.ObjectId	`json:"id,omitempty" bson:"_id,omitempty"`
		Username	string	`json:"username,omitempty" bson:"username,omitempty"`
		Firstname	string	`json:"firstname,omitempty" bson:"firstname,omitempty"`
		Lastname	string `json:"lastname,omitempty" bson:"lastname,omitempty"`
		Password	string	`json:"password,omitempty" bson:"password,omitempty"`
		Image	string	`json:"image,omitempty" bson:"image,omitempty"`
	}
)

func (u *User)IsNotDuplicate() bool {
	err := helper.UsersCollection.Find(bson.M{"username": u.Username}).One(u)
	if err != nil {
		return true
	}
	return false
}

func (u *User)CheckLogin() (*User,error){
	err := helper.UsersCollection.Find(bson.M{"username": u.Username,"password": u.Password}).One(&u)
	if err != nil {
		return nil, err 
	}
	return u,nil
}

func (u *User)SaveUserToDB() error {
	err := helper.UsersCollection.Insert(&u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User)ReadUsersFromDB() ([]User , error ){
	result := []User{}
	err := helper.UsersCollection.Find(nil).All(&result)
	if err != nil {
		return nil,err
	}
	return result, nil
}

func (u *User)ReadUsersByIDFromDB() (*User, error) {
	err := helper.UsersCollection.Find(bson.M{"_id": u.Id}).One(&u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User)UpdateUserToDB() bool{
	err := helper.UsersCollection.UpdateId(u.Id,u)
	if err != nil {
		return false
	}
	return true
}

func (u *User)DeleteUserByIDFromDB() bool {
	err := helper.UsersCollection.RemoveId(u.Id)
	if err != nil {
		return false
	}
	return true
}
