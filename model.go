package main

import "time"

type Person struct {
	id int
	first_name string
	last_name string
	birth time.Time
	gender int
	misc *string
}

type Relation struct {
	id int
	src int
	dest int
	relationType int
}
