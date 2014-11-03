package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"log"
	"net"
	//"os/exec"
	//"time"
)

const (
	CONN_HOST = ""
	CONN_PORT = "7474"
	CONN_TYPE = "tcp"
)

func main() {

	var debug = flag.Bool("d", false, "debug - whether to print copious what-i-m-doing messages")
	flag.Parse()

	var keywurds Keywurds

	var b Bot
	b.Name = "AIGOR"
	b.Mood = 100
	b.Debug = *debug
	b.innit(&keywurds)

	db, err := sql.Open("postgres", "user=thingie dbname=thingie sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Thing{}, "thing").SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		panic("Error Listening:" + err.Error())
	}
	conns := clientConns(l)
	for {
		go b.talkPerson(<-conns, keywurds, dbmap)
	}
}

func clientConns(l net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		fmt.Println("\n\n*********\nYar, Chatbot Listening on " + CONN_HOST + ":" + CONN_PORT)
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Err accepting: ", err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, conn.LocalAddr(), conn.RemoteAddr())
			ch <- conn
		}
	}()
	return ch
}
