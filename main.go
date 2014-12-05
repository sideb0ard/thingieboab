package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"strings"
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
	done := make(chan bool)

	mypid := os.Getpid()
	fmt.Println("my pid is", strconv.Itoa(mypid))

	go b.yinmind(neurons)
	go b.yangmind(neurons)

	//go b.netchat(dbmap)

	initThought := Thought{"what am i, who am i, when am i", 100}
	neurons <- initThought

	<-done

}

func (b Bot) yinmind(neurons chan Thought) {
	for t := range neurons {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("yin", t)
		fmt.Println("yin", t.Wurds)
		wurds := strings.Split(t.Wurds, " ")
		replymsg := []string{}
		for i, w := range wurds {
			if i%7 != 0 {
				replymsg = append(replymsg, w)
			} else {
				fmt.Println("BOOP")
			}
		}
		fmt.Println("len", len(replymsg))
		fmt.Println("join", strings.Join(replymsg, " "))
		reply := Thought{strings.Join(replymsg, " "), 99}
		neurons <- reply
	}
}

func (b Bot) yangmind(neurons chan Thought) {
	for t := range neurons {
		fmt.Println("yang", t)
		fmt.Println("yangw", t.Wurds)
		wurds := strings.Split(t.Wurds, " ")
		replymsg := []string{}
		for i, w := range wurds {
			if i%3 == 0 {
				for i := 0; i < 2; i++ {
					replymsg = append(replymsg, w)
				}
			} else {
				replymsg = append(replymsg, w)
			}
		}
		reply := Thought{strings.Join(replymsg, " "), 99}
		neurons <- reply
	}
}
