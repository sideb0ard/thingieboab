package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	"io/ioutil"
	"net/http"
	//"log"
	"net"
	//"regexp"
	//"strconv"
	"strings"
	//"time"
)

func (b Bot) innit(keywurds *Keywurds) {
	file, err := ioutil.ReadFile("/var/server/bobbybot.json")
	if err != nil {
		file, _ = ioutil.ReadFile("./bobbybot.json")
	}

	err = json.Unmarshal(file, &keywurds)
	if err != nil {
		fmt.Println("ERRRRRRR:", err)
	}

}

func (b Bot) talkPerson(conn net.Conn, keywurds Keywurds, dbmap *gorp.DbMap) {
	bu := bufio.NewReader(conn)
	p := &Thing{}

	conn.Write([]byte(">** HUllo. I am " + b.Name + "\n"))

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
	for {
		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
			break
		}
		reply := b.transform(string(line), keywurds)
		conn.Write([]byte("\n>" + reply + "\n"))
	}
}

func (b Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}

func (b Bot) transform(s string, keywurds Keywurds) string {
	url := "http://localhost:5000/"
	var jsonnn = "{\"wurds\": \"" + strings.TrimSpace(s) + "\"}"
	fmt.Println(jsonnn)
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

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return string(body)

}
