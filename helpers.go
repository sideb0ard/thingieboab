package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	//fmt.Println("MAX IS ", max, " MIN IS ", min)
	return rand.Intn(max-min) + min
	//return 7
}

func tokenizer(bigstringy string) []string {
	// returns a slice from sentence of all smaller phrases
	var tokenz []string
	var innertoke func(str string)
	innertoke = func(str string) {
		tokenz = append(tokenz, str)
		var stringies = strings.Split(str, " ")
		if len(stringies) == 1 {
			return
		}
		innertoke(strings.Join(stringies[0:len(stringies)-1], " "))
	}
	var stringies = strings.Split(bigstringy, " ")
	for i := range stringies {
		innertoke(strings.Join(stringies[i:], " "))
	}
	return tokenz
}

func spaceify(dashy string) string {
	spxex, _ := regexp.Compile("-")
	spaceyString := strings.ToLower(spxex.ReplaceAllString(dashy, " "))
	return spaceyString
}
func dashify(spacey string) string {
	spxex, _ := regexp.Compile(" ")
	dashyString := strings.ToLower(spxex.ReplaceAllString(spacey, "-"))
	return dashyString
}
func saveKnowledge(thing string, meaning string) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	key := strings.ToLower(thing)
	_, err = c.Do("SET", key, meaning)
	if err != nil {
		fmt.Println(err)
	}
}
func getKnowledge(q string) string {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	rkey := "aigor:memory:" + strings.ToLower(q)
	r, err := redis.String(c.Do("GET", rkey))
	if err != nil {
		fmt.Println(err)
	}
	if len(r) > 0 {
		return spaceify(r)
	} else {
		return ""
	}
}
func getKeys(q string) []string {
	var keys []string
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	keys, err = redis.Strings(c.Do("KEYS", strings.ToLower(q)+"*"))
	if err != nil {
		fmt.Println(err)
	}
	return keys
}
func getValue(k string) string {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	v, _ := redis.Bytes(c.Do("GET", strings.ToLower(k)))
	return string(v)
}

//func (b Bot) understand(sentence string, person_name string) (string, string) {
//	// first check for pronouns, then local redis, then internet dictionary
//	var subject string
//	var action string
//	var wurds = strings.Split(sentence, " ")
//	// combine pronouns and known names into a pool of likely subject candidates.
//	//subjects := pronouns
//	subjects[strings.ToLower(b.Name)] = 1
//	subjects[strings.ToLower(person_name)] = 1
//	for w := range wurds {
//		_, ok := subjects[strings.ToLower(wurds[w])]
//		if ok {
//			r := regexp.MustCompile(`(?i)\b(` + wurds[w] + `)\b(.*)`)
//			matches := r.FindAllStringSubmatch(sentence, -1)
//			subject = matches[0][1]
//			action = matches[0][2]
//
//		}
//	}
//	return subject, action
//}
