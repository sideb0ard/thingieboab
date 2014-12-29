package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	//"strings"
	"time"
)

func main() {

	var debug = flag.Bool("d", false, "debug - whether to print copious what-i-m-doing messages")
	flag.Parse()

	var b Bot
	b.Name = "RED"
	b.Mood = 100
	b.Debug = *debug

	startTime := time.Now().Unix()
	fmt.Println("came online at ", startTime)
	db, err := sql.Open("postgres", "user=thingie dbname=thingie sslmode=disable")
	err = db.Ping()
	if err != nil {
		log.Fatal("Is postgres running? ", err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(ThingType{}, "thingType").SetKeys(true, "Id")
	dbmap.AddTableWithName(Thing{}, "thing").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	neurons := make(chan Thought)
	memory := make(chan Event)
	done := make(chan bool)

	mypid := os.Getpid()
	fmt.Println("my pid is", strconv.Itoa(mypid))

	go b.yinmind(neurons, memory)
	go b.yangmind(neurons, memory)

	go b.storeMem(memory)

	//go b.netchat(dbmap)

	//initThought := Thought{"what am i, who am i, when am i", 100}
	initThought := Thought{"I am not a number! I am a free man!", 100}
	neurons <- initThought

	<-done

}

func (b Bot) storeMem(memory chan Event) {
	for m := range memory {
		fmt.Printf("MEMORY:: when: %s // what: %s // ", m.When, m.What)
		b.CurrentStimulii = append(b.CurrentStimulii, m)
		fmt.Println("Memories LEN: ", len(b.CurrentStimulii))
	}
}

func (b Bot) yinmind(neurons chan Thought, memory chan Event) {
	for t := range neurons {
		time.Sleep(time.Duration(random(0, 1000)) * time.Millisecond)
		//fmt.Println("yin", t.Wurds)

		wurds := strings.Split(t.Wurds, " ")
		//var replywurds string
		for i := range wurds {
			j := rand.Intn(i + 1)
			wurds[i], wurds[j] = wurds[j], wurds[i]
		}

		mem := Event{time.Now(), "neuron received", "myself"}
		memory <- mem

		reply := Thought{strings.Join(wurds, " "), 99}
		fmt.Println("YIN SENDING REPLY", reply)
		neurons <- reply
	}
}

func (b Bot) yangmind(neurons chan Thought, memory chan Event) {
	for t := range neurons {
		time.Sleep(time.Duration(random(0, 1000)) * time.Millisecond)
		//fmt.Println("yangw", t.Wurds)

		wurds := strings.Split(t.Wurds, " ")
		//var replywurds string
		for i := range wurds {
			j := rand.Intn(i + 1)
			wurds[i], wurds[j] = wurds[j], wurds[i]
		}

		mem := Event{time.Now(), "neuron received", "myself"}
		memory <- mem

		reply := Thought{strings.Join(wurds, " "), 99}
		fmt.Println("YaANG INSENDING REPLY", reply)
		neurons <- reply
	}
}
