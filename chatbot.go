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
	"strings"
	"time"
)

var pronouns map[string]int

func (b Bot) innit() {
	pronouns = b.convertRedisKey("pronoun")
	file, _ := ioutil.ReadFile("./reasmb_rules.json")
	//reass_file, e := os.Open("./reasmb_rules.json")

	var kw Keywurds
	err := json.Unmarshal(file, &kw)
	if err != nil {
		fmt.Println("ERRRRRRR:", err)
	}

	//fmt.Println(kw.Keywords)
	for blah := range kw.Keywords {
		fmt.Println(blah)
		fmt.Println(kw.Keywords[blah].Score)
	}

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

func (b Bot) talkperson(bored_chan chan bool, listen_chan chan string, mood_chan chan int, neurons_chan chan Thought) {
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
			reply := transform(line)
			reply = b.procsz(reply, "post")
			fmt.Println(reply)
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

func transform(s string) string {
	//rank := -2
	//reasmb := ""
	re := regexp.MustCompile(`[?!,]`)
	s = re.ReplaceAllString(s, ".")
	rebut := regexp.MustCompile(`but`)
	s = rebut.ReplaceAllString(s, ".")
	sparts := strings.Split(s, ".")
	fmt.Println(sparts)
	return s
}
