package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var people []Person
var db *sql.DB
var err error

func main() {
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/", getIndex)
	router.POST("/person", addPerson)
	router.LoadHTMLFiles("index.html")
	router.StaticFile("styles.css", "styles.css")
	router.StaticFile("favicon.svg", "favicon.svg")
	router.Run(":8080")
}

func personToView(p Person) PersonView {
	return PersonView{
		ID:       p.id,
		FullName: p.firstName + " " + p.lastName,
		Age:      calculateAge(p.birth),
		Gender:   p.gender,
	}
}

func getIndex(c *gin.Context) {
	people = getPeople(db)
	p := people[rand.Intn(len(people))]
	for _, v := range people {
		if v.firstName == "hakima" {
			p = v
			break
		}
	}
	view := personToView(p)

	// Populate parents
	parents := getParents(db, p.id)
	for _, parent := range parents {
		view.Parents = append(view.Parents, personToView(parent))
	}

	// Populate siblings
	siblings := getSiblings(db, p.id)
	for _, sib := range siblings {
		view.Siblings = append(view.Siblings, personToView(sib))
	}

	// Populate siblings
	children := getChildren(db, p.id)
	for _, chi := range children {
		view.Children = append(view.Children, personToView(chi))
	}

	// Populate partners
	partners := getPartners(db, p.id)
	for _, partner := range partners {
		view.Partners = append(view.Partners, personToView(partner))
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"person":   view,
		"childnum": len(children) - 1,
	})
}

func addPerson(c *gin.Context) {
	p := Person{}
	p.firstName = c.PostForm("first_name")
	p.lastName = c.PostForm("last_name")
	p.gender = c.PostForm("gender")
	birth := c.PostForm("birth")
	p.birth, _ = time.Parse("2006-01-02", birth)
	death := c.PostForm("death")
	if death != "" {
		d, _ := time.Parse("2006-01-02", death)
		p.death = &d
	}

	newId, err := savePerson(db, p)
	if err != nil {
		log.Fatal("failed to add person: ", err)
	}
	fatherId := c.PostForm("fatherId")
	motherId := c.PostForm("motherId")
	relation := Relation{}
	relation.father, _ = strconv.Atoi(fatherId)
	relation.mother, _ = strconv.Atoi(motherId)
	relation.person = newId
	err = saveRelation(db, relation)
	if err != nil {
		log.Fatal("failed to add person: ", err)
	}

	people = getPeople(db)
	var views []PersonView
	for _, p := range people {
		views = append(views, PersonView{
			FullName: p.firstName + " " + p.lastName,
			Age:      calculateAge(p.birth),
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
