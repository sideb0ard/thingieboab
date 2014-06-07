package main

import (
	"bufio"
	"fmt"
	"io"
	//"io/ioutil"
	//"log"
	//"encoding/json"
	//"math/rand"
	"os"
	"regexp"
	"strings"
)

var pronouns map[string]int

//var coordinators map[string]int

func (b Bot) innit() {
	pronouns = b.convertRedisKey("pronoun")
	//fmt.Println(pronouns)
}

func listen() string {
	bio := bufio.NewReader(os.Stdin)
	line, err := bio.ReadBytes('\n')
	if err == io.EOF {
		os.Exit(0)
	}
	if err != nil {
		panic(err)
	}
	sline := strings.TrimRight(string(line), "\n")
	return sline
}

func (b Bot) talkPerson() {
	p := Person{}

	fmt.Printf("Hullo. I am " + b.name + "\n")

	if len(p.name) == 0 {
		fmt.Println("What is your name?")
		p.name = listen()
		fmt.Printf("YOU ARE %v\nWhat you want to talk about?\n", p.name)
	}

	for {
		line := listen()
		//fmt.Println("ORIG LINE IS ", line)
		line = b.procsz(line, "pre")
		//fmt.Println("PRE-PROCESSED LINE IS ", line)

		//fmt.Println("POST-PROCESSED LINE IS ", line)
		subject, action := understand(line)
		//fmt.Println("SUBJECT:" + subject)

		//subject = b.procsz(subject, "post")

		if len(subject) > 0 {
			you := regexp.MustCompile(`(?i)\bi\b`)
			me := regexp.MustCompile(`(?i)\byou\b`)
			action = b.procsz(action, "post")
			if me.MatchString(subject) {
				fmt.Println("ITs ME!", b.name+", I"+action)
			} else if you.MatchString(subject) {
				fmt.Println("ITs YOU!", p.name+", YOU"+action)
			} else {
				fmt.Println("Oh yeah, " + subject + ". Yeah, " + action)
			}
		} else {
			bangqregex := regexp.MustCompile(`[!?]`)
			line = bangqregex.ReplaceAllString(line, "")
			reply := getValue("aigor:memory:" + line)
			if len(reply) > 0 {
				fmt.Println(reply)
			} else {
				fmt.Printf("Sorry, i don't know what %v means - can you tell me?\n", line)
				explanation := listen()
				fmt.Printf("Thanks, so \"%v\" means \"%v\" - got it (i think!!)\n", line, explanation)
				saveKnowledge(b.name+":memory:"+line, explanation)
			}
		}

	}
}

func (b Bot) think() {
	question := "What am I?"
	answer := "I think, therefore I am."
	saveKnowledge(question, answer)
}

func (b Bot) dream() {
	fmt.Println("electric sheepzzzzzzz")
}
func (b Bot) convertRedisKey(pre string) map[string]int {
	prefix := strings.ToLower(b.name) + ":" + pre + ":"
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
	prefix := strings.ToLower(b.name) + ":" + stage + ":"
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
