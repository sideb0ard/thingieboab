package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//var pronouns map[string]int

func (b Bot) innit(keywurds *Keywurds) {
	//pronouns = b.convertRedisKey("pronoun")
	//if b.Debug {
	//	fmt.Println("PRONOUNSZZ:", pronouns)
	//}
	file, err := ioutil.ReadFile("/var/server/bobbybot.json")
	if err != nil {
		file, _ =  ioutil.ReadFile("./bobbybot.json")
	}

	err = json.Unmarshal(file, &keywurds)
	if err != nil {
		fmt.Println("ERRRRRRR:", err)
	}

}


func (b Bot) uppermind(mood_chan chan int, neurons_chan chan Thought) {
	for t := range neurons_chan {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(t)
		reply := Thought{"UPPER MIND SHIZZLE IT", <-mood_chan}
		neurons_chan <- reply
	}
}

func (b Bot) lowermind(mood_chan chan int, neurons_chan chan Thought) {
	for t := range neurons_chan {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(t)
		//wurds := search5(t.wurds)
		reply := Thought{"LOWER MINDS WURDZZZ", <-mood_chan}
		neurons_chan <- reply
	}
}

func (b Bot) moody(mood_chan chan int) {
	for {
		mood_chan <- rand.Intn(100)
	}
}

func (b Bot) think(bored_chan chan bool, mood_chan chan int, neurons_chan chan Thought) {
	randy := rand.Intn(2)
	switch {
	case randy == 0:
		bored_chan <- true
	case randy == 1:
		neurons_chan <- Thought{"**random Th0ught **", <-mood_chan}
	}
}

//func (b Bot) talkPerson(bored_chan chan bool, listen_chan chan string, mood_chan chan int, neurons_chan chan Thought, keywurds Keywurds) {
func (b Bot) talkPerson(conn net.Conn, keywurds Keywurds) {
	bu := bufio.NewReader(conn)
	p := Thing{}

	//fmt.Printf(">**Hullo. I am " + b.Name + "\n")
	conn.Write([]byte(">** HUllo. I am " + b.Name + "\n"))

	if len(p.Name) == 0 {
		//fmt.Println(">What is your name?")
		conn.Write([]byte(">What is your name*?\n"))
		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
		}
		//p.Name = <-listen_chan
		p.Name = strings.TrimSpace(string(line))
		//fmt.Printf("\n>Pleased to meet ya %v\n>Wha's gon' on?\n\n", p.Name)
		conn.Write([]byte("\n>Please to meet ya " + p.Name + ". Wha's up?\n"))
	}
	for {
		line, err := bu.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errzzz reading :", err.Error())
			break
		}

		//select {
		//case line, _ := <-listen_chan:
		//line = b.procsz(line, "pre")
		reply := b.transform(string(line), keywurds)
		conn.Write([]byte("\n>" + reply + "\n"))
		//fmt.Println(">", reply, "\n")
		//subject, action := b.understand(line, p.Name)
		//if len(subject) > 0 {
		//	action = b.procsz(action, "post")
		//	if regexp.MustCompile(`(?i)\byou\b`).MatchString(subject) || regexp.MustCompile(`(?i)\b`+b.Name+`\b`).MatchString(subject) {
		//		fmt.Println("ITs ME! ", b.Name)
		//	} else if regexp.MustCompile(`(?i)\bi\b`).MatchString(subject) || regexp.MustCompile(`(?i)\b`+p.Name+`\b`).MatchString(subject) {
		//		fmt.Println("ITs YOU!", p.Name+", YOU"+action)
		//	} else {
		//		fmt.Println("Oh yeah, " + subject + ". Yeah, " + action)
		//	}
		//} else {
		//	bangqregex := regexp.MustCompile(`[!?]`)
		//	line = bangqregex.ReplaceAllString(line, "")
		//	reply := getValue("aigor:memory:" + line)
		//	if len(reply) > 0 {
		//		fmt.Println(reply)
		//	} else {
		//		fmt.Printf("Sorry, i don't know what %v means - can you tell me?\n", line)
		//		explanation, _ := <-listen_chan
		//		fmt.Printf("Thanks, so \"%v\" means \"%v\" - got it (i think!!)\n", line, explanation)
		//		saveKnowledge(b.Name+":memory:"+line, explanation)
		//	}
		//}

		//case Thought, _ := <-neurons_chan:
		//	fmt.Println("HERES A Thought...", Thought)
		//case _, _ = <-bored_chan:
		//	fmt.Println("**bzzzt** getting bored here, challenge me, bro.. **8zz8**")
		//}
	}
}

func (b Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}

//func (b Bot) convertRedisKey(pre string) map[string]int {
//	prefix := strings.ToLower(b.Name) + ":" + pre + ":"
//	fullkeys := getKeys(prefix)
//	re, _ := regexp.Compile(prefix + `(.*)`)
//	keys := make(map[string]int)
//	for i := range fullkeys {
//		keys[string(re.FindStringSubmatch(fullkeys[i])[1])] = 1
//	}
//	return keys
//}

func (b Bot) procsz(s string, stage string) string {
	// pre- or post- process a string and return updated string
	prefix := strings.ToLower(b.Name) + ":" + stage + ":"
	//fullkeys := getKeys(prefix)
	//re, _ := regexp.Compile(prefix + `(.*)`)
	keys := make(map[string]int)
	//for i := range fullkeys {
	//	keys[string(re.FindStringSubmatch(fullkeys[i])[1])] = 1
	//}
	//if b.Debug {
	//	fmt.Println("IN PROCSZ with", s)
	//	fmt.Println("PREFIX", prefix)
	//	fmt.Println("FULLKEYS:", fullkeys)
	//	fmt.Println("KEYS:", keys)
	//}
	var wurds = strings.Split(s, " ")
	for w := range wurds {
		if b.Debug {
			fmt.Println("Looking for ", wurds[w], " in ", keys)
		}
		_, ok := keys[wurds[w]]
		if ok {
			r := regexp.MustCompile(`\b` + wurds[w] + `\b`)
			s = r.ReplaceAllString(s, getValue(prefix+wurds[w]))
		}
	}
	if b.Debug {
		fmt.Println("RETURNINF FROM PROCSZ with", s)
	}
	return s
}

func (b Bot) transform(s string, keywurds Keywurds) string {
	if b.Debug {
		fmt.Println("I GOTZ ", s)
	}
	s = b.procsz(s, "pre")
	if b.Debug {
		fmt.Println("AN NOW I GOTZ ", s)
	}
	score := -2
	reasmb := ""
	re := regexp.MustCompile(`[?!,]`)
	s = re.ReplaceAllString(s, ".")
	rebut := regexp.MustCompile(`but`)
	s = rebut.ReplaceAllString(s, ".")
	sparts := strings.Split(s, ".")

	reassemblrrr := func(kw string, spart int) string {
		for d := range keywurds.Keywords[kw].Decomp {
			// match the asterix in the decomp rules
			dre := regexp.MustCompile(`(?i)\s*\*\s*`)
			// change them to word boundary capture groups
			nre := regexp.MustCompile(`(?i)` + dre.ReplaceAllString(d, "\\b(.*)\\b"))
			decomp_matches := nre.FindAllStringSubmatch(sparts[spart], -1)
			if b.Debug {
				fmt.Println("D iz: ", d)
				fmt.Println("DECOMPZ: ", decomp_matches)
			}
			if len(decomp_matches) > 0 {

				// grab a random reassembly rule
				randy := random(0, len(keywurds.Keywords[kw].Decomp[d]))
				reasmb = keywurds.Keywords[kw].Decomp[d][randy]

				// if its a goto, return the new keyword so it can reassembled
				goto_re := regexp.MustCompile(`^goto\s(\w*).*`)
				if goto_re.FindString(reasmb) != "" {
					nkw := goto_re.FindStringSubmatch(reasmb)
					return (nkw[1])
				}
				exec_re := regexp.MustCompile(`^ish\s(\w*).*`)
				if exec_re.FindString(reasmb) != "" {
					reasmb = time.Now().Format(time.RFC1123)
				}

				for j := 1; j < len(decomp_matches[0]); j++ {
					if decomp_matches[0][j] == "" {
						continue
					}
					decomp_matches[0][j] = b.procsz(decomp_matches[0][j], "post")
					decomp_re := regexp.MustCompile(`\(` + strconv.Itoa(j) + `\)`)
					reasmb = decomp_re.ReplaceAllString(reasmb, decomp_matches[0][j])
				}
			}
			// TODO Match synonyms - @
		}
		return ""
	}

	for spart := range sparts {
		for kw := range keywurds.Keywords {
			if b.Debug {
				fmt.Println("LOOKING FOR ", kw, " IN ", sparts[spart])
			}
			if regexp.MustCompile(`(?i)\b`+kw+`\b`).MatchString(sparts[spart]) && score < keywurds.Keywords[kw].Score {
				score = keywurds.Keywords[kw].Score
				rerun := reassemblrrr(kw, spart)
				if rerun != "" {
					reassemblrrr(rerun, spart)
				}
			}
		}
	}

	if reasmb == "" {
		randy := random(0, len(keywurds.Keywords["xnone"].Decomp["*"]))
		reasmb = keywurds.Keywords["xnone"].Decomp["*"][randy]
	}
	spx := regexp.MustCompile(`\s+`)
	reasmb = spx.ReplaceAllString(reasmb, " ")
	q_re := regexp.MustCompile(`\s\?`)
	reasmb = q_re.ReplaceAllString(reasmb, "?")
	return reasmb

}
