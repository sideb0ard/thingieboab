package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	CONN_HOST = ""
	CONN_PORT = "7474"
	CONN_TYPE = "tcp"
)

func (b *Bot) netchat(memory chan Event) {

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		panic("Error Listening:" + err.Error())
	}
	conns := clientConns(l)
	for {
		go b.talkHandler(<-conns, memory)
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

func (b *Bot) talkHandler(conn net.Conn, memory chan Event) {
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
	p.Name = nom

	memory <- Event{time.Now(), "met", p}

	known, err := b.retrieveMemory(p)
	if err != nil {
		fmt.Println("EERZZZ", err.Error())
	}
	fmt.Println("KNOWN::", known)
	if len(known) > 0 {
		conn.Write([]byte("\n>hey " + nom + ", I know you!. How are ya?\n"))
	} else {
		conn.Write([]byte("\n>Please to meet ya " + nom + ". Wha's up?\n"))
	}

	//questionasked := false
	//var currentsubject string
	memory <- Event{time.Now(), "starting chat with", p}

	for {

		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
			break
		}

		memory <- Event{time.Now(), "asked question by", p}

		//	sline := strings.TrimSpace(string(line))
		//	if questionasked {
		//		fmt.Println("QUESTION HAS BEEN ASKED - current subject is", currentsubject)
		//		t := &Thing{}
		//		//retrieve <- *t
		//		//known := <-remember
		//		//if known == true {
		//		//	t.Properties = sline
		//		//	store <- *t
		//		//}
		//		questionasked = false
		//	}
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
				//t := &Thing{}
				//retrieve <- *t
				//known := <-remember
				//if known == false {
				//	t.Name = s
				//	t.Type = "person"
				//	store <- *t
				//	conn.Write([]byte("\n>Who or what is a " + s + "?\n"))
				//	currentsubject = s
				//	questionasked = true
				//} else {
				//	fmt.Println("I know it!")
				//	conn.Write([]byte("\n>Oh yeah, " + s + " I know of this\n"))
				//}
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
