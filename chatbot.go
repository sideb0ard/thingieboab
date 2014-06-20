package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	//"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var pronouns map[string]int

func (b Bot) innit(keywurds *Keywurds) {
	pronouns = b.convertRedisKey("pronoun")
	file, _ := ioutil.ReadFile("./bobbybot.json")
	//reass_file, e := os.Open("./reasmb_rules.json")

	err := json.Unmarshal(file, &keywurds)
	if err != nil {
		fmt.Println("ERRRRRRR:", err)
	}

	//fmt.Println(kw.Keywords)
	//m := new(Dispatch)
	//var m interface{}
	//var jsontype jsonobject
	//json.Unmarshal(file, &jsontype)
	//fmt.Printf("Results: %v\n", jsontype)
}

func (b Bot) listen(listen_chan chan string) {
	bio := bufio.NewReader(os.Stdin)
	for {
		line, err := bio.ReadBytes('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		sline := strings.TrimRight(string(line), "\n")
		listen_chan <- sline
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

func (b Bot) talkPerson(bored_chan chan bool, listen_chan chan string, mood_chan chan int, neurons_chan chan Thought, keywurds Keywurds) {
	p := Person{}

	fmt.Printf("Hullo. I am " + b.Name + "\n")

	if len(p.Name) == 0 {
		fmt.Println("What is your name?")
		p.Name = <-listen_chan
		fmt.Printf("Pleased to meet ya %v\nWha's gon' on?\n", p.Name)
	}
	for {
		select {
		case line, _ := <-listen_chan:
			line = b.procsz(line, "pre")
			reply := transform(line, keywurds)
			reply = b.procsz(reply, "post")
			//fmt.Println("MAIN REPLY", reply)
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

		case Thought, _ := <-neurons_chan:
			fmt.Println("HERES A Thought...", Thought)
		case _, _ = <-bored_chan:
			fmt.Println("**bzzzt** getting bored here, challenge me, bro.. **8zz8**")
		}
	}
}

func (b Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}
func (b Bot) convertRedisKey(pre string) map[string]int {
	prefix := strings.ToLower(b.Name) + ":" + pre + ":"
	fullkeys := getKeys(prefix)
	re, _ := regexp.Compile(prefix + `(.*)`)
	keys := make(map[string]int)
	for i := range fullkeys {
		keys[string(re.FindStringSubmatch(fullkeys[i])[1])] = 1
	}
	return keys
}

func (b Bot) procsz(s string, stage string) string {
	// pre- or post- process a string and return updated string
	prefix := strings.ToLower(b.Name) + ":" + stage + ":"
	fullkeys := getKeys(prefix)
	re, _ := regexp.Compile(prefix + `(.*)`)
	keys := make(map[string]int)
	for i := range fullkeys {
		keys[string(re.FindStringSubmatch(fullkeys[i])[1])] = 1
	}
	var wurds = strings.Split(s, " ")
	for w := range wurds {
		_, ok := keys[wurds[w]]
		if ok {
			r := regexp.MustCompile(`\b` + wurds[w] + `\b`)
			s = r.ReplaceAllString(s, getValue(prefix+wurds[w]))

		}
	}
	return s
}

func transform(s string, keywurds Keywurds) string {
	score := -2
	//reasmb := ""
	re := regexp.MustCompile(`[?!,]`)
	s = re.ReplaceAllString(s, ".")
	rebut := regexp.MustCompile(`but`)
	s = rebut.ReplaceAllString(s, ".")
	sparts := strings.Split(s, ".")
	for spart := range sparts {
		//fmt.Println("SPART: ", sparts[spart])
		for kw := range keywurds.Keywords {
			// fmt.Println("KW: ", kw, "SCORE IS ", score, "KEYWURD SCORE IS ", keywurds.Keywords[kw].Score)
			if regexp.MustCompile(`(?i)\b`+kw+`\b`).MatchString(sparts[spart]) && score < keywurds.Keywords[kw].Score {
				//fmt.Println("KEYWORD MATCH (and lower score) : ", keywurds.Keywords[kw].Decomp, spart, score)
				score = keywurds.Keywords[kw].Score
				for d := range keywurds.Keywords[kw].Decomp {
					dre := regexp.MustCompile(`(?i)\s*\*\s*`)
					nre := regexp.MustCompile(`(?i)` + dre.ReplaceAllString(d, "\\b(.*)\\b"))
					decomp_matches := nre.FindAllStringSubmatch(sparts[spart], -1)
					if len(decomp_matches) > 0 {
						//fmt.Println("DECOMP MATCH! : ", decomp_matches)

						//fmt.Println("WUP _ MATCHED!")
						//fmt.Println("DECOMP: ", d)
						//fmt.Println("DECOMP_MATCHES: ", decomp_matches)
						//fmt.Println("RECOMPREPLY: ", nre)
						////fmt.Println(len(decomp_matches))
						//randy := random(1, len(decomp_matches[0]))
						//fmt.Println("LENNY: ", len(decomp_matches[0]))
						//fmt.Println("RANDY: ", randy)
						//fmt.Println("RANDREPLY:", decomp_matches[0][randy])
						//fmt.Println("**D**:", keywurds.Keywords[kw].Decomp[d])
						//fmt.Println("**LEN D**:", len(keywurds.Keywords[kw].Decomp[d]))
						randy := random(0, len(keywurds.Keywords[kw].Decomp[d]))
						reply := keywurds.Keywords[kw].Decomp[d][randy]
						for j := range decomp_matches[0] {
							decomp_re := regexp.MustCompile(`\(` + strconv.Itoa(j) + `\)`)
							//fmt.Println("J: ", strconv.Itoa(j))
							//fmt.Println("ITEM: ", decomp_matches[0][j])
							//fmt.Println("DECOMP_RE: ", decomp_re)
							reply = decomp_re.ReplaceAllString(reply, decomp_matches[0][j])
							spx := regexp.MustCompile(`\s+`)
							reply = spx.ReplaceAllString(reply, " ")
						}
						fmt.Println(reply)
						return (reply)
					}
					// TODO Match synonyms - @
				}
			}
		}
	}
	//	fmt.Println(blah)
	//	//	//fmt.Println(kw.Keywords[blah].Score)
	//}

	return s
}
