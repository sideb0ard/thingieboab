package main

type Thought struct {
	Wurds string
	Mood  int
}

type Bot struct {
	Name  string
	Mood  int
	Debug bool
}

type Thing struct {
	Id        int64
	Name      string
	Mood      int
	ThingType string
	//Properties    []interface{}
	//Relationships []interface{}
	//Memories      []interface{}
}

type TransformReply struct {
	Breakdown string
}
