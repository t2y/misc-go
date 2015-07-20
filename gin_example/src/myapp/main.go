package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var Logger = logrus.New()

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserInfo struct {
	users []User `json:"users"`
}

func SetProperties(c *gin.Context) {
	c.Set("myapp", "something")
	c.Next()
}

func getDB(c *gin.Context) *sql.DB {
	db, ok := c.Get("DB")
	if ok {
		return db.(*sql.DB)
	}
	return nil
}

func ListUsers(c *gin.Context) {
	rows := QuerySql(getDB(c), "select id, name, age from users")
	users := make([]User, 0, 10)
	for rows.Next() {
		var id int
		var name string
		var age int
		err := rows.Scan(&id, &name, &age)
		if err == nil {
			users = append(users, User{Id: id, Name: name, Age: age})
		}
	}
	defer rows.Close()

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	paramName := c.Param("name")
	row := QueryRowSql(getDB(c), "select id, name, age from users where name = ?", paramName)
	if row == nil {
		Logger.Error("not found")
		return
	}

	var id int
	var name string
	var age int
	err := row.Scan(&id, &name, &age)
	if err == nil {
		c.JSON(http.StatusOK, User{Id: id, Name: name, Age: age})
	} else {
		Logger.Errorf(`%v`, err)
	}
}

func RegisterUser(c *gin.Context) {
	type Param struct {
		Id  int `json:"id"`
		Age int `json:"age"`
	}
	var p Param
	if c.BindJSON(&p) == nil {
		InsertUser(getDB(c), p.Id, c.Param("name"), p.Age)
		msg := fmt.Sprintf("inserted id: %d", p.Id)
		c.JSON(http.StatusOK, gin.H{"status": msg})
	}
}

func UpdateUserInfo(c *gin.Context) {
	type Param struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var p Param
	if c.BindJSON(&p) == nil {
		oldName := c.Param("name")
		UpdateUser(getDB(c), oldName, p.Name, p.Age)
	}
}

func DeleteUserInfo(c *gin.Context) {
	DeleteUser(getDB(c), c.Param("name"))
}

func main() {
	Logger.Formatter = new(logrus.JSONFormatter)

	// init
	db := InitDatabase("./myapp.db")
	defer db.Close()

	r := gin.Default()

	// middleware
	r.Use(func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	})
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(SetProperties)

	// handlers
	r.GET("/users/", ListUsers)
	r.GET("/users/:name", GetUser)
	r.POST("/users/:name", RegisterUser)
	r.PUT("/users/:name", UpdateUserInfo)
	r.DELETE("/users/:name", DeleteUserInfo)

	r.Run(":8000")
}
