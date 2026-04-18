package main

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var people []Person
var db *sql.DB
var err error

func main() {
	db, err = sql.Open("sqlite3", "file:database.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/", getIndex)
	router.GET("/add", getAdd)
	router.POST("/person", addPerson)
	router.DELETE("/person", deletePerson)
	router.POST("/relation", addRelation)
	router.LoadHTMLFiles("index.html", "add.html")
	router.StaticFile("styles.css", "styles.css")
	router.StaticFile("favicon.svg", "favicon.svg")
	router.Run(":8181")
}

func personToView(p Person) PersonView {
	return PersonView{
		ID:       p.id,
		FullName: p.firstName + " " + p.lastName,
		Age:      calculateAge(p.birth),
		Gender:   p.gender,
	}
}

func getAdd(c *gin.Context) {
	people = getPeople(db)
	if len(people) == 0 {
		c.String(http.StatusNotFound, "no people found")
		return
	}

	var allPeopleViews []PersonView
	for _, person := range people {
		allPeopleViews = append(allPeopleViews, personToView(person))
	}
	c.HTML(http.StatusOK, "add.html", gin.H{
		"people": allPeopleViews,
	})
}

func getIndex(c *gin.Context) {
	people = getPeople(db)
	if len(people) == 0 {
		c.String(http.StatusNotFound, "no people found")
		return
	}

	var allPeopleViews []PersonView
	for _, person := range people {
		allPeopleViews = append(allPeopleViews, personToView(person))
	}

	selectedPersonID, convErr := strconv.Atoi(c.Query("person_id"))
	p := people[rand.Intn(len(people))]
	if convErr == nil {
		for _, person := range people {
			if person.id == selectedPersonID {
				p = person
				break
			}
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
	childNum := 0
	if len(children) > 0 {
		childNum = len(children) - 1
	}

	// Populate partners
	partners := getPartners(db, p.id)
	for _, partner := range partners {
		view.Partners = append(view.Partners, personToView(partner))
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"person":   view,
		"childnum": childNum,
		"people":   allPeopleViews,
	})
}

func addPerson(c *gin.Context) {
	p := Person{}
	p.firstName = c.PostForm("first_name")
	p.lastName = c.PostForm("last_name")
	p.gender = c.PostForm("gender")
	if p.gender != "F" && p.gender != "M" {
		c.String(http.StatusBadRequest, "invalid gender")
		return
	}
	birth := c.PostForm("birth")
	p.birth, _ = time.Parse("2006-01-02", birth)
	death := c.PostForm("death")
	if death != "" {
		d, _ := time.Parse("2006-01-02", death)
		p.death = &d
	}

	newId, err := savePerson(db, p)
	if err != nil {
		log.Printf("failed to add person: %v", err)
		c.String(http.StatusInternalServerError, "failed to add person")
		return
	}

	fatherRaw := strings.TrimSpace(c.PostForm("fatherId"))
	motherRaw := strings.TrimSpace(c.PostForm("motherId"))
	fatherID, convErr := parseOptionalID(fatherRaw)
	if convErr != nil {
		c.String(http.StatusBadRequest, "invalid father id")
		return
	}
	motherID, convErr := parseOptionalID(motherRaw)
	if convErr != nil {
		c.String(http.StatusBadRequest, "invalid mother id")
		return
	}
	if fatherID != nil || motherID != nil {
		relation := Relation{}
		relation.person = newId
		relation.father = fatherID
		relation.mother = motherID
		err = saveRelation(db, relation)
		if err != nil {
			log.Printf("failed to add relation for person %d: %v", newId, err)
			c.String(http.StatusInternalServerError, "failed to add relation")
			return
		}
	}

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}

func addRelation(c *gin.Context) {
	personRaw := strings.TrimSpace(c.PostForm("personId"))
	fatherRaw := strings.TrimSpace(c.PostForm("fatherId"))
	motherRaw := strings.TrimSpace(c.PostForm("motherId"))

	if personRaw == "" {
		c.String(http.StatusBadRequest, "person is required")
		return
	}

	relation := Relation{}
	var convErr error
	relation.person, convErr = strconv.Atoi(personRaw)
	if convErr != nil {
		c.String(http.StatusBadRequest, "invalid person id")
		return
	}
	relation.father, convErr = parseOptionalID(fatherRaw)
	if convErr != nil {
		c.String(http.StatusBadRequest, "invalid father id")
		return
	}
	relation.mother, convErr = parseOptionalID(motherRaw)
	if convErr != nil {
		c.String(http.StatusBadRequest, "invalid mother id")
		return
	}
	if relation.father == nil && relation.mother == nil {
		c.String(http.StatusBadRequest, "father or mother is required")
		return
	}

	err := saveRelation(db, relation)
	if err != nil {
		log.Printf("failed to add relation: %v", err)
		c.String(http.StatusInternalServerError, "failed to add relation")
		return
	}

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}

func deletePerson(c *gin.Context) {
	personRaw := strings.TrimSpace(c.Query("personId"))
	if personRaw == "" {
		personRaw = strings.TrimSpace(c.PostForm("personId"))
	}
	if personRaw == "" {
		personRaw = readDeleteField(c, "personId")
	}
	if personRaw == "" {
		c.String(http.StatusBadRequest, "person is required")
		return
	}

	personID, convErr := strconv.Atoi(personRaw)
	if convErr != nil {
		c.String(http.StatusBadRequest, "invalid person id")
		return
	}

	err := deletePersonByID(db, personID)
	if errors.Is(err, sql.ErrNoRows) {
		c.String(http.StatusNotFound, "person not found")
		return
	}
	if err != nil {
		log.Printf("failed to delete person: %v", err)
		c.String(http.StatusInternalServerError, "failed to delete person")
		return
	}

	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}

func readDeleteField(c *gin.Context, field string) string {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		return ""
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(values.Get(field))
}

func parseOptionalID(raw string) (*int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func calculateAge(birth time.Time) int {
	now := time.Now()
	age := now.Year() - birth.Year()
	if now.Month() < birth.Month() || (now.Month() == birth.Month() && now.Day() < birth.Day()) {
		age--
	}
	return age
}
