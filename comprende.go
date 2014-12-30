package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (b *Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}

func (b *Bot) comprende(wurds []string) (s []string, o []string) {
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

func (b *Bot) transform(s string) []byte {
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
