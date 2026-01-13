package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var people []Person

func main() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	people = getPersons(db)
	router := gin.Default()
	router.GET("/", getIndex)
	router.LoadHTMLFiles("index.html")
	router.StaticFile("styles.css", "styles.css")
	router.StaticFile("config.js", "config.js")
	router.StaticFile("favicon.svg", "favicon.svg")
	router.Run(":8080")
}

func getIndex(c *gin.Context) {
	person := people[0]
	c.HTML(http.StatusOK, "index.html", gin.H{
		"fullName": person.first_name + " " + person.last_name,
		"age":      calculateAge(person.birth),
	})
}

func calculateAge(birth time.Time) int {
	now := time.Now()
	age := now.Year() - birth.Year()
	if now.Month() < birth.Month() || (now.Month() == birth.Month() && now.Day() < birth.Day()) {
		age--
	}
	return age
}

func getPersons(db *sql.DB) (people []Person) {
	stmnt := "select * from person"
	rows, err := db.Query(stmnt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var person Person
		err := rows.Scan(
			&person.id,
			&person.first_name,
			&person.last_name,
			&person.birth,
			&person.gender,
			&person.misc)
		if err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
		i++
	}
	return
}
