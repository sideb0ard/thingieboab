package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
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

func main() {

	var debug = flag.Bool("d", false, "debug - whether to print copious what-i-m-doing messages")
	flag.Parse()

	var b Bot
	b.Name = "AIGOR"
	b.Mood = 100
	b.Debug = *debug

	db, err := sql.Open("postgres", "user=thingie dbname=thingie sslmode=disable")
	err = db.Ping()
	if err != nil {
		log.Fatal("Is postgres running? ", err)
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
		go b.talkHandler(<-conns, dbmap)
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

func (b Bot) talkHandler(conn net.Conn, dbmap *gorp.DbMap) {
	bu := bufio.NewReader(conn)
	//p := &Thing{}

	conn.Write([]byte(">** HUllo. I am " + b.Name + "\n"))

	// first run - get some initial infozzzz...
	//if len(p.Name) == 0 {
	//	conn.Write([]byte(">What is your name*?\n"))
	//	line, err := bu.ReadBytes('\n')
	//	if err != nil {
	//		fmt.Println("Errzzz reading :", err.Error())
	//	}
	//	nom := strings.TrimSpace(string(line))
	//	err = dbmap.SelectOne(p, "select * from thing where name=$1", nom)
	//	if err != nil {
	//		fmt.Println(nom, "doesnt exist - now creating..")
	//		conn.Write([]byte("\n>Please to meet ya " + nom + ". Wha's up?\n"))
	//		p.Name = nom
	//		err = dbmap.Insert(p)
	//		if err != nil {
	//			fmt.Println("ERrrr inserting to DB:", err.Error())
	//		}
	//	} else {
	//		fmt.Println("woop! exists")
	//		conn.Write([]byte("\n>hey " + nom + ", I know you!. How are ya?\n"))
	//	}
	//}
	for {
		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
			break
		}
		sentence := b.transform(string(line))
		var s TransformReply
		err = json.Unmarshal(sentence, &s)
		if err != nil {
			fmt.Println("Ooft, json not unmarshalling...")
		}
		wurds := strings.Split(s.Breakdown, " ")
		comprende(wurds)

		//fmt.Println("SUBJ/OBJ: ", subject, " + ", object)
		//fmt.Println("WURDS: ", wurds)
		//resj := regexp.MustCompile()
		conn.Write([]byte("\n>" + strings.Join(wurds, " // ") + "\n"))
	}
}

func (b Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}

func comprende(wurds []string) (s string, o string) {
	var sentence []Wurd
	re, _ := regexp.Compile("([a-zA-Z0-9'-]+)/([A-Z0-9-]+)/([A-Z0-9-]+)/([A-Z0-9-]+)/([A-Z0-9-]+)/([A-Z0-9-]+)/([a-zA-Z0-9'-]+)")
	for w := range wurds {
		if re.MatchString(wurds[w]) {
			tw := re.FindStringSubmatch(wurds[w])
			id, _ := strconv.Atoi(tw[4])
			anch, _ := strconv.Atoi(tw[6])
			wobj := Wurd{
				Name:   tw[1],
				Tag:    tw[2],
				Chunk:  tw[3],
				Id:     id,
				Role:   tw[5],
				Anchor: anch,
				Lemma:  tw[7],
			}
			sentence = append(sentence, wobj)
		}
	}
	for i := range sentence {
		fmt.Println(sentence[i].String())
	}
	return "subject", "object"
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
