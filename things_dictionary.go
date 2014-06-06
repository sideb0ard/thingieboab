package main

type Person struct {
	name string
	mood int
}

type Thought struct {
	wurds string
	mood  int
}

type Bot struct {
	name string
	mood int
}

type Thing struct {
	name          string
	thingType     string
	properties    []interface{}
	relationships []interface{}
	memories      []interface{}
}