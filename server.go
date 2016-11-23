package main

import (
	"net/http"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"./models"
)

var (
	mongoSession *mgo.Session
)


func getSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")
    if err != nil {
        panic(err)
    }
    return session
}

func saveUser(c echo.Context) error{
	user := new(models.User)
	user.Id = bson.NewObjectId()
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "please check input")
	}
	if user.IsNotDuplicate(mongoSession) && !(user.Username == "") {
		user.SaveUserToDB(mongoSession)
		return c.JSON(http.StatusCreated, "create user success")
	}
	return c.JSON(http.StatusConflict, "user empty or duplicate")
}

func getUser(c echo.Context) error {
	user := models.User{}
	results := []models.User{}
	results, _ = user.ReadUsersFromDB(mongoSession)
	return c.JSON(http.StatusOK, results)
}

func getUserByID(c echo.Context) error {
	user := models.User{}
	id := c.Param("id")
	user.Id = bson.ObjectIdHex(id)
	results, _ := user.ReadUsersByIDFromDB(mongoSession)
	return c.JSON(http.StatusOK, results)
}

func updateUser(c echo.Context) error {
	user := new(models.User)
	id := c.Param("id")
	user.Id = bson.ObjectIdHex(id)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "please check input")
	}
	if user.IsNotDuplicate(mongoSession) {
		if user.UpdateUserToDB(mongoSession) {
			return c.JSON(http.StatusOK, "update user completed")
		}
	}
	return c.JSON(http.StatusBadRequest, "can not up to date user") 
}

func deleteUser(c echo.Context) error {
	user := models.User{}
	id := c.Param("id")
	user.Id = bson.ObjectIdHex(id)
	if user.DeleteUserByIDFromDB(mongoSession) {
		return c.JSON(http.StatusOK, "delete user completed")
	}
	return c.JSON(http.StatusBadRequest, "not found user or can not to delete")
}

func init() {
	mongoSession = getSession()
	mongoSession.SetMode(mgo.Monotonic, true)
}

func main() {
	e := echo.New()

	e.Static("/images", "images")

	e.POST("/users", saveUser)
	e.GET("/users", getUser)
	e.GET("/users/:id", getUserByID)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
