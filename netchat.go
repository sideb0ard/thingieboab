package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

const (
	CONN_HOST = ""
	CONN_PORT = "7474"
	CONN_TYPE = "tcp"
)

func (b *Bot) netchat(dbmap *gorp.DbMap, memory chan Event) {

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		panic("Error Listening:" + err.Error())
	}
	conns := clientConns(l)
	for {
		go b.talkHandler(<-conns, dbmap, memory)
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

func (b *Bot) talkHandler(conn net.Conn, dbmap *gorp.DbMap, memory chan Event) {
	bu := bufio.NewReader(conn)
	p := &Thing{}

	conn.Write([]byte(">** Hello. I am " + b.Name + "\n"))

	memory <- Event{time.Now(), "said hello", p}

	// first run - get some initial infozzzz...
	conn.Write([]byte(">What is your name*?\n"))
	line, err := bu.ReadBytes('\n')
	if err != nil {
		fmt.Println("Errzzz reading :", err.Error())
	}
	nom := strings.TrimSpace(string(line))

	memory <- Event{time.Now(), "met", p}

	err = dbmap.SelectOne(p, "select * from thing where name=$1", nom)
	if err != nil {
		fmt.Println(nom, "doesnt exist - now creating..")
		conn.Write([]byte("\n>Please to meet ya " + nom + ". Wha's up?\n"))
		p.Name = nom
		err = dbmap.Insert(p)
		if err != nil {
			fmt.Println("ERrrr inserting to DB:", err.Error())
		}
	} else {
		fmt.Println("woop! exists")
		conn.Write([]byte("\n>hey " + nom + ", I know you!. How are ya?\n"))
	}

	questionasked := false
	var currentsubject string
	memory <- Event{time.Now(), "starting chat with", p}

	for {

		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
			break
		}

		memory <- Event{time.Now(), "asked question by", p}

		sline := strings.TrimSpace(string(line))
		if questionasked {
			fmt.Println("QUESTION HAS BEEN ASKED - current subject is", currentsubject)
			t := &Thing{}
			err := dbmap.SelectOne(t, "select * from thing where name=$1", currentsubject)
			fmt.Println("Ok, t is:", t)
			if err == nil {
				t.Properties = sline
				count, err := dbmap.Update(t)
				if err != nil {
					fmt.Println("Beep! error inserting to DB:", err.Error())
				}
				fmt.Println("Updated ", count, " rows..")
			} else {
				fmt.Println("errrr:", err.Error())
			}
			questionasked = false
		}
		sentence := b.transform(string(line))
		fmt.Println("SENTENCE", string(sentence))
		var s TransformReply
		err = json.Unmarshal(sentence, &s)
		if err != nil {
			fmt.Println("Ooft, json not unmarshalling...")
		}
		wurds := strings.Split(s.Breakdown, " ")
		subjects, objects := b.comprende(wurds)

		fmt.Println("STIMULUS:::ZZZ::: ", b.CurrentStimulii)

		fmt.Println("SUBJ/OBJ: ", subjects, " + ", objects)
		if len(subjects) > 0 {
			for _, s := range subjects {
				fmt.Println("Oh, yeah, ", s)
				t := &Thing{}
				err := dbmap.SelectOne(t, "select * from thing where name=$1", s)
				if err != nil {
					t.Name = s
					t.Type = "person"
					err = dbmap.Insert(t)
					if err != nil {
						fmt.Println("Beep! error inserting to DB:", err.Error())
					}
					conn.Write([]byte("\n>Who or what is a " + s + "?\n"))
					currentsubject = s
					questionasked = true
				} else {
					fmt.Println("I know it!")
					conn.Write([]byte("\n>Oh yeah, " + s + " I know of this\n"))
				}
			}

		} else {
			conn.Write([]byte("\n>I'm sorry, that does not compute. I am but a simple A.I.\n"))
		}
		//fmt.Println("WURDS: ", wurds)
		//resj := regexp.MustCompile()
		//conn.Write([]byte("\n>" + strings.Join(wurds, " // ") + "\n"))
		//conn.Write([]byte("\n> Ha! you're talking about " + subject + " and " + object + "\n"))
	}
}
