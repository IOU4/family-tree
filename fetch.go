package main

import (
	"database/sql"
	"log"
)

func scanPerson(rows *sql.Rows, person *Person) error {
	var death sql.NullTime
	err := rows.Scan(
		&person.id,
		&person.firstName,
		&person.lastName,
		&person.birth,
		&death,
		&person.gender,
		&person.misc,
	)
	if err != nil {
		return err
	}
	if death.Valid {
		person.death = &death.Time
	} else {
		person.death = nil
	}
	return nil
}

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
		err := scanPerson(rows, &person)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
		i++
	}
	return
}

func getParents(db *sql.DB, personID int) []Person {
	query := `
		SELECT p.* FROM person p
		WHERE p.id = (select father from relation where person = ?)
		OR p.id = (select mother from relation where person = ?)
	`
	rows, err := db.Query(query, personID, personID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var parents []Person
	for rows.Next() {
		var person Person
		err := scanPerson(rows, &person)
		if err != nil {
			log.Fatal(err)
		}
		parents = append(parents, person)
	}

	return parents
}

func getSiblings(db *sql.DB, personID int) []Person {
	query := `
		select p.* from person p join relation r on r.person = p.id
		where (r.father = (select father from relation where person = ?)
		or r.mother = (select mother from relation where person = ?))
		and p.id != ?
	`

	rows, err := db.Query(query, personID, personID, personID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err := scanPerson(rows, &person)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
	}
	return people
}

func getChildren(db *sql.DB, personID int) []Person {
	query := `
		select p.* from person p join relation r on r.person = p.id
		where r.father = ? or r.mother = ?
	`

	rows, err := db.Query(query, personID, personID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err := scanPerson(rows, &person)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
	}
	return people
}

func getPartners(db *sql.DB, personID int) []Person {
	query := `
		select p.* from person p
		where p.id in (
			select case
				when r.father = ? then r.mother
				when r.mother = ? then r.father
			end
			from relation r
			where r.father = ? or r.mother = ?
		)
	`

	rows, err := db.Query(query, personID, personID, personID, personID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var person Person
		err := scanPerson(rows, &person)
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
	}
	return people
}

func savePerson(db *sql.DB, person Person) (int, error) {
	result, err := db.Exec(
		"INSERT INTO person (first_name, last_name, birth, death, gender) VALUES (?, ?, ?, ?, ?)",
		person.firstName, person.lastName, person.birth, person.death, person.gender)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func saveRelation(db *sql.DB, relation Relation) error {
	_, err := db.Exec("INSERT INTO relation (person, father, mother) values (?, ?, ?)",
		relation.person, relation.father, relation.mother)
	return err
}
