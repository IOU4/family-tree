package main

import (
	"database/sql"
	"log"
)

func getPeople(db *sql.DB) (people []Person) {
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
			&person.firstName,
			&person.lastName,
			&person.birth,
			&person.gender,
			&person.misc)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
		i++
	}
	return
}

func savePerson(db *sql.DB, person Person) error {
	_, err := db.Exec(
		"INSERT INTO person (first_name, last_name, birthday, gender) VALUES (?, ?, ?, ?)",
		person.firstName, person.lastName, person.birth, person.gender)
	return err
}

func saveRelation(db *sql.DB, relation Relation) error {
	_, err := db.Exec("INSERT INTO relation (p1, p2, kinship) values (?, ?, ?)",
		relation.p1, relation.p2, relation.kinship)
	return err
}

func getPersonByID(db *sql.DB, id int) (Person, error) {
	var person Person
	err := db.QueryRow("SELECT * FROM person WHERE id = ?", id).Scan(
		&person.id, &person.firstName, &person.lastName,
		&person.birth, &person.gender, &person.misc)
	return person, err
}

func getRelatedPeople(db *sql.DB, personID int, kinship string) []Person {
	// Find people related to personID with the given kinship type.
	// Check both directions of the relation (p1->p2 and p2->p1).
	query := `
		SELECT p.* FROM person p
		JOIN relation r ON (r.p1 = p.id OR r.p2 = p.id)
		WHERE (r.p1 = ? OR r.p2 = ?) AND r.kinship = ? AND p.id != ?
	`
	rows, err := db.Query(query, personID, personID, kinship, personID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.id, &person.firstName, &person.lastName,
			&person.birth, &person.gender, &person.misc)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
	}
	return people
}

func getParents(db *sql.DB, personID int) []Person {
	query := `
		SELECT p.* FROM person p
		JOIN relation r ON r.p1 = p.id
		WHERE r.p2 = ? AND r.kinship = 'P'
	`
	rows, err := db.Query(query, personID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var parents []Person
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.id, &person.firstName, &person.lastName,
			&person.birth, &person.gender, &person.misc)
		if err != nil {
			log.Fatal(err)
		}
		parents = append(parents, person)
	}

	// If we found one parent, also find their partner
	if len(parents) == 1 {
		partners := getRelatedPeople(db, parents[0].id, "H")
		parents = append(parents, partners...)
	}

	return parents
}
