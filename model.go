package main

import "time"

type Person struct {
	id        int
	firstName string
	lastName  string
	birth     time.Time
	death     *time.Time
	gender    string
	misc      *string
}

type Relation struct {
	person int
	father *int
	mother *int
}

type PersonView struct {
	ID       int
	FullName string
	Age      int
	Gender   string
	Siblings []PersonView
	Partners []PersonView
	Parents  []PersonView
	Children []PersonView
}
