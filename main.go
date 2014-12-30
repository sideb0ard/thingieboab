package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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

	//go b.yinMind(neurons, memory)
	//go b.yangMind(neurons, memory)

	go b.storeMind(memory)

	go b.netchat(dbmap, memory)

	//initThought := Thought{"what am i, who am i, when am i", 100}
	initThought := Thought{"I am not a number! I am a free man!", 100}
	neurons <- initThought

	<-done

}
