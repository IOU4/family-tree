package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var people []Person
var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	people = getPersons(db)
	router := gin.Default()
	router.GET("/", getIndex)
	router.POST("/person", addPerson)
	router.LoadHTMLFiles("index.html")
	router.StaticFile("styles.css", "styles.css")
	router.StaticFile("favicon.svg", "favicon.svg")
	router.Run(":8080")
}

func getIndex(c *gin.Context) {
	var views []PersonView
	for _, p := range people {
		views = append(views, PersonView{
			FullName: p.FirstName + " " + p.LastName,
			Age:      calculateAge(p.Birth),
		})
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"people": views,
	})
}

func addPerson(c *gin.Context) {
	firstName := c.PostForm("first_name")
	lastName := c.PostForm("last_name")
	birthday := c.PostForm("birthday")
	genderStr := c.PostForm("gender")

	gender, err := strconv.Atoi(genderStr)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid gender")
		return
	}

	_, err = db.Exec(
		"INSERT INTO person (first_name, last_name, birthday, gender) VALUES (?, ?, ?, ?)",
		firstName, lastName, birthday, gender,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to add person")
		return
	}

	// refresh in-memory list
	people = getPersons(db)

	// return updated grid for htmx swap
	var views []PersonView
	for _, p := range people {
		views = append(views, PersonView{
			FullName: p.FirstName + " " + p.LastName,
			Age:      calculateAge(p.Birth),
		})
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"people": views,
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
			&person.ID,
			&person.FirstName,
			&person.LastName,
			&person.Birth,
			&person.Gender,
			&person.Misc)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
		i++
	}
	return
}
