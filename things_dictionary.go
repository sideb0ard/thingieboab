package main

import (
	"fmt"
	"time"

	"github.com/coopernurse/gorp"
)

type Thought struct {
	Wurds string
	Mood  int
}

type Event struct {
	When time.Time
	What string
	Who  *Thing
	//Associations []Thing
	//Source       Thing
}

func (e Event) String() string {
	return fmt.Sprintf("<<%s %s %s>>", e.When, e.What, e.Who.Name)
}

type Bot struct {
	Thing
	CurrentStimulii []Event
	Db              *gorp.DbMap
}

type ThingType struct {
	Id   int64
	Type string
	//RelationType []string
	//Properties []Property
	Properties string
	//Memories     []string
}

type Thing struct {
	Id    int64
	Name  string
	Mood  int
	Debug bool
	ThingType
	//Properties    []interface{}
	//Relationships []interface{}
	//Memories      []interface{}
}

type Property struct {
	Id     int64
	Detail string
}

type RelationType struct {
	Id        int64
	Name      string
	Source    ThingType
	Recipient ThingType
}

type Relation struct {
	Id        int64
	Name      string
	Source    Thing
	Recipient Thing
	//Type      []RelationType
	Type string
}

type TransformReply struct {
	Breakdown string
}

type Wurd struct {
	// from - http://www.clips.ua.ac.be/pages/MBSP
	Word     string
	POS      string // Part-of-speech
	Chunk    string
	PNP      int // prepositional noun phrase
	Relation string
	//PNP    string
	Anchor int
	Lemma  string
}

func (w Wurd) String() string {
	return fmt.Sprintf("<<%q %q %q %d %q %d %q>>", w.Word, w.POS, w.Chunk, w.PNP, w.Relation, w.Anchor, w.Lemma)
}
