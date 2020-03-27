package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"log"
	"strconv"
)

type User struct {
	Id   int    `form:"id" binding:"required" db:"id" json:"id"`
	Name string `form:"name" db:"name" json:"name"`
}

var (
	DBCon *gorp.DbMap
)

func DbOpen() *gorp.DbMap {
	db, err := sql.Open("mysql", "not_needed_test:none")
	if err != nil {
		log.Fatalln("sql.Open failed", err)
	}

	DBCon := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	return DBCon
}

func initDb(dbcon *gorp.DbMap) *gorp.DbMap {
	dbcon.AddTableWithName(User{}, "user").SetKeys(true, "Id")
	err := dbcon.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Create tables failed", err)
	}

	return dbcon
}

func GetMainEngine(dbcon *gorp.DbMap) *gin.Engine {
	DBCon = dbcon
	r := gin.Default()
	v1 := r.Group("api/v1")
	{
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
	}
	v1.Use()

	return r
}

func main() {
	dbcon := DbOpen()
	dbcon = initDb(dbcon)
	GetMainEngine(dbcon).Run(":8080")
}

func GetUsers(c *gin.Context) {
	var users []User
	_, err := DBCon.Select(&users, "SELECT * FROM user")

	if err == nil {
		c.JSON(200, users)
	} else {
		c.JSON(404, gin.H{"error": "no users found!"})
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	err := DBCon.SelectOne(&user, "SELECT * FROM user WHERE id=? LIMIT 1", id)

	if err == nil {
		userId, _ := strconv.Atoi(id)

		content := &User{
			Id:   userId,
			Name: user.Name,
		}
		c.JSON(200, content)
	} else {
		log.Fatal(err)
		c.JSON(404, gin.H{"error": "user not found"})
	}
}
