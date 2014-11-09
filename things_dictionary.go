package main

import (
	"fmt"
)

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

type Wurd struct {
	Name  string
	Tag   string
	Chunk string
	Id    int
	Role  string
	//PNP    string
	Anchor int
	Lemma  string
}

func (w Wurd) String() string {
	return fmt.Sprintf("<<%q %q %q %d %q %d %q>>", w.Name, w.Tag, w.Chunk, w.Id, w.Role, w.Anchor, w.Lemma)
}
