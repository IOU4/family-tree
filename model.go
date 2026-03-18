package main

import "time"

type Person struct {
	id        int
	firstName string
	lastName  string
	birth     time.Time
	gender    string
	misc      *string
}

type Relation struct {
	id        int
	p1				int
	p2        int
	kinship		string
}

type PersonView struct {
	ID				int
	FullName	string
	Age				int
	Gender		string
	Siblings	[]PersonView
	Parents 	[]PersonView
}
