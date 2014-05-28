package main

import (
	"bufio"
	"fmt"
	//"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"math/rand"
	"os"
	"regexp"
	//"strconv"
	"strings"
)

func (b Bot) innit() {
	file, err := os.Open("language.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewReader(file)
	for {
		line, err := scanner.ReadBytes('\n')
		if err == io.EOF {
			//os.Exit(0)
			break
		}
		if err != nil {
			panic(err)
		}
		if match, _ := regexp.Match(`^#`, line); match == true {
			continue
		}
		sline := strings.Split(strings.TrimRight(string(line), "\n"), "|")
		entryType, keyWord, replacement := sline[0], sline[1], strings.Join(sline[2:], " ")
		storageKey := b.name + ":" + entryType + ":" + keyWord
		//fmt.Println(storageKey, replacement)
		saveKnowledge(storageKey, replacement)
	}
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
		//fmt.Println("LINE:", line)
		//line = b.procsz(line, "pre")
		//fmt.Println("PRELINE:", line)

		reply := getKnowledge(line)
		fmt.Println("REPLY:", reply)
		reply = b.procsz(reply, "post")
		fmt.Println("POSTREPLY:", reply)

		if len(reply) == 0 {
			b.mood -= 10
			//fmt.Println("LAME - no REPLY!")
			fmt.Printf("Sorry, i don't know what %v means - can you tell me?\n", line)
			explanation := listen()
			fmt.Printf("Thanks, so \"%v\" means \"%v\" - got it (i think!!)\n", line, explanation)
			saveKnowledge(b.name+":memory:"+line, explanation)
		}
		sayName := rand.Intn(4)
		if sayName == 0 {
			b.mood += 10
			reply = p.name + ", " + reply
		}
		b.mood += 10
		fmt.Println(reply)
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

func (b Bot) procsz(s string, stage string) string {
	// pre- or post- process a string and return updated string
	prefix := strings.ToLower(b.name) + ":" + stage + ":"
	fullkeys := getKeys(prefix)
	re, _ := regexp.Compile(prefix + `(.*)`)
	//fmt.Println(re)
	keys := make(map[string]int)
	for i := range fullkeys {
		fmt.Println("YA", string(re.FindStringSubmatch(fullkeys[i])[1]))

		keys[string(re.FindStringSubmatch(fullkeys[i])[1])] = 1
		//fmt.Println(getValue(string(fullkeys[i])))
	}
	var wurds = strings.Split(s, " ")
	fmt.Println("INIT WURDS: ", s)
	for w := range wurds {
		_, ok := keys[wurds[w]]
		if ok {
			r := regexp.MustCompile(wurds[w])
			s = r.ReplaceAllString(s, getValue(prefix+wurds[w]))

		}
	}
	fmt.Println("RETURN WURDS: ", s)

	return s
}
