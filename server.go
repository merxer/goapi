package main

import (
	"fmt"
	"os"
	"time"
	"strconv"
	"net/http"
	b64 "encoding/base64"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"./models"
	"./helper"
)

const (
	SERVERNAME="128.199.130.61:80"
	//SERVERNAME="localhost:1323"
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
	if user.IsNotDuplicate() && !(user.Username == "") {
		user.SaveUserToDB()
		return c.JSON(http.StatusCreated, "create user success")
	}
	return c.JSON(http.StatusConflict, "user empty or duplicate")
}

func getUser(c echo.Context) error {
	user := models.User{}
	results := []models.User{}
	results, _ = user.ReadUsersFromDB()
	return c.JSON(http.StatusOK, results)
}

func getUserByID(c echo.Context) error {
	user := models.User{}
	id := c.Param("id")
	user.Id = bson.ObjectIdHex(id)
	results, _ := user.ReadUsersByIDFromDB()
	return c.JSON(http.StatusOK, results)
}

func updateUser(c echo.Context) error {
	user := new(models.User)
	id := c.Param("id")
	user.Id = bson.ObjectIdHex(id)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "please check input")
	}
	fmt.Printf("%#v",user)
	if user.IsNotDuplicate() {
		if user.UpdateUserToDB() {
			return c.JSON(http.StatusOK, "update user completed")
		}
	}
	return c.JSON(http.StatusBadRequest, "can not up to date user")
}

func deleteUser(c echo.Context) error {
	user := models.User{}
	id := c.Param("id")
	user.Id = bson.ObjectIdHex(id)
	if user.DeleteUserByIDFromDB() {
		return c.JSON(http.StatusOK, "delete user completed")
	}
	return c.JSON(http.StatusBadRequest, "not found user or can not to delete")
}

func uploadImage(c echo.Context) error {
	image := new(models.Image)
	image.Id = bson.NewObjectId()
	image.Name =  strconv.FormatInt(time.Now().UTC().UnixNano(),10)
	if err := c.Bind(image); err != nil {
		return c.JSON(http.StatusBadRequest, "please check image input")
	}
	sDec, _ := b64.StdEncoding.DecodeString(image.Base64String)
	file, err := os.Create("./images/"+image.Name+".png")
	if err != nil {
		panic("please check images path")
	}
	defer file.Close()
	_, err = file.Write(sDec)
	if err != nil {
		panic(err)
	}
    result := new(models.Image)
	result.Url = "http://"+ SERVERNAME+ "/images/"+ image.Name + ".png"
	return c.JSON(http.StatusOK, result)
}

func login(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, "please check input")
	}

	result, err  := user.CheckLogin()
	if err != nil {
		return echo.ErrUnauthorized
	}
	return c.JSON(http.StatusOK, result)
}

func init() {
	mongoSession := getSession()
	mongoSession.SetMode(mgo.Monotonic, true)
	helper.MongoSession = mongoSession
	helper.UsersCollection = mongoSession.DB("mdb").C("users")
}

func main() {
	defer helper.MongoSession.Close()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))


	e.Static("/images", "images")

	e.POST("/users", saveUser)
	e.GET("/users", getUser)
	e.GET("/users/:id", getUserByID)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	e.POST("/upload/images", uploadImage)
	e.POST("/login", login)

	e.Logger.Fatal(e.Start(SERVERNAME))
}
