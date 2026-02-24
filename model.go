package main

import "time"

type Person struct {
	ID        int
	FirstName string
	LastName  string
	Birth     time.Time
	Gender    int
	Misc      *string
}

type Relation struct {
	id           int
	src          int
	dest         int
	relationType int
}

type PersonView struct {
	FullName string
	Age      int
}
