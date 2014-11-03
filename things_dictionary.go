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

type Keywurds struct {
	Keywords map[string]Decomp
}

type Decomp struct {
	Score  int
	Decomp map[string][]string
}
