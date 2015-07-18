package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/goji/glogrus"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
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

func SetProperties(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Pass data through the environment
		c.Env["myapp"] = "something"
		// Fully control how the next layer is called
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getDB(c web.C) *sql.DB {
	db, ok := c.Env["DB"].(*sql.DB)
	if ok {
		return db
	}
	return nil
}

func ListUsers(c web.C, w http.ResponseWriter, r *http.Request) {
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

	bytes, _ := json.Marshal(users)
	w.Write(bytes)
}

func GetUser(c web.C, w http.ResponseWriter, r *http.Request) {
	paramName := c.URLParams["name"]
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
		user := User{Id: id, Name: name, Age: age}
		bytes, _ := json.Marshal(user)
		w.Write(bytes)
	} else {
		Logger.Errorf(`%v`, err)
	}
}

func RegisterUser(c web.C, w http.ResponseWriter, r *http.Request) {
	type Param struct {
		Id  int `json:"id"`
		Age int `json:"age"`
	}
	var p Param
	json.NewDecoder(r.Body).Decode(&p)
	InsertUser(getDB(c), p.Id, c.URLParams["name"], p.Age)
	msg := fmt.Sprintf("inserted id: %d", p.Id)
	w.Write([]byte(msg))
}

func UpdateUserInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	type Param struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var p Param
	json.NewDecoder(r.Body).Decode(&p)
	oldName := c.URLParams["name"]
	UpdateUser(getDB(c), oldName, p.Name, p.Age)
}

func DeleteUserInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	DeleteUser(getDB(c), c.URLParams["name"])
}

func main() {
	Logger.Formatter = new(logrus.JSONFormatter)

	// init
	db := InitDatabase("./myapp.db")
	defer db.Close()

	// middleware
	goji.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["DB"] = db
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)

	})
	goji.Use(glogrus.NewGlogrus(Logger, "myapp"))
	goji.Use(middleware.NoCache)
	goji.Use(SetProperties)

	// handlers
	goji.Get("/users/", ListUsers)
	goji.Get(regexp.MustCompile(`/users/(?P<name>\w+)$`), GetUser)
	goji.Post(regexp.MustCompile(`/users/(?P<name>\w+)$`), RegisterUser)
	goji.Put("/users/:name", UpdateUserInfo)
	goji.Delete("/users/:name", DeleteUserInfo)

	goji.Serve()
}
