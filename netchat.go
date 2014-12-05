package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	CONN_HOST = ""
	CONN_PORT = "7474"
	CONN_TYPE = "tcp"
)

func (b Bot) netchat(dbmap *gorp.DbMap) {

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		panic("Error Listening:" + err.Error())
	}
	conns := clientConns(l)
	for {
		//go b.talkHandler(<-conns, dbmap)
		go b.initialGreeting(<-conns, dbmap)
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

func (b Bot) initialGreeting(conn net.Conn, dbmap *gorp.DbMap) {
	bu := bufio.NewReader(conn)

	p := &Thing{}

	conn.Write([]byte(">** Hello. I am " + b.Name + "\n"))

	// first run - get some initial infozzzz...
	if len(p.Name) == 0 {
		conn.Write([]byte(">What is your name*?\n"))
		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
		}
		nom := strings.TrimSpace(string(line))
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
	}
	b.talkHandler(conn, dbmap)
}

func (b Bot) talkHandler(conn net.Conn, dbmap *gorp.DbMap) {
	bu := bufio.NewReader(conn)
	questionasked := false
	var currentsubject string
	for {

		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
			break
		}
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

func (b Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}

func (b Bot) comprende(wurds []string) (s []string, o []string) {
	var sentence []Wurd
	re, _ := regexp.Compile(`([\$a-zA-Z0-9'-]+)/([\$A-Z0-9-]+)/([\$A-Z0-9-]+)/([\$A-Z0-9-]+)/([\$A-Z0-9-]+)/([\$A-Z0-9-]+)/([\$a-zA-Z0-9'-]+)`)
	for w := range wurds {
		if re.MatchString(wurds[w]) {
			tw := re.FindStringSubmatch(wurds[w])
			id, _ := strconv.Atoi(tw[4])
			anch, _ := strconv.Atoi(tw[6])
			wobj := Wurd{
				Word:     tw[1],
				POS:      tw[2],
				Chunk:    tw[3],
				PNP:      id,
				Relation: tw[5],
				Anchor:   anch,
				Lemma:    tw[7],
			}
			sentence = append(sentence, wobj)
		}
	}
	if b.Debug {
		for i := range sentence {
			fmt.Println(sentence[i].String())
		}
	}
	var sj, oj []string
	sre, _ := regexp.Compile("SBJ")
	ore, _ := regexp.Compile("OBJ")
	for w := range sentence {
		if sre.MatchString(sentence[w].Relation) {
			sj = append(sj, sentence[w].Word)
		} else if ore.MatchString(sentence[w].Relation) {
			oj = append(oj, sentence[w].Word)
		}
	}
	return sj, oj
}

func (b Bot) transform(s string) []byte {
	url := "http://localhost:5000/"
	var jsonnn = "{\"wurds\": \"" + strings.TrimSpace(s) + "\"}"
	var jsonStr = []byte(jsonnn)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Sideb0ard-Service", "v1")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if b.Debug {
		fmt.Println(jsonnn)
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
	}
	return body

}
