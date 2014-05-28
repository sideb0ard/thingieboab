package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"strings"
)

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
	//fmt.Println("SPACEME CALLED! Sending back: " + spacyString)
	return spaceyString
}
func dashify(spacey string) string {
	spxex, _ := regexp.Compile(" ")
	dashyString := strings.ToLower(spxex.ReplaceAllString(spacey, "-"))
	//fmt.Println("DASHME CALLED! Sending back: " + dashyString)
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

	rkey := "aigor:memory:" + q
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
	fmt.Println("WOOP, IN GETKEYS")
	var keys []string
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	//key := "aigor:memory:" + q
	fmt.Println("Q IS :" + q)
	keys, err = redis.Strings(c.Do("KEYS", strings.ToLower(q)+"*"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(keys)
	return keys
}
func getValue(k string) string {
	//fmt.Println("WOOP, IN GET VALUE!")
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	//fmt.Println("Key IS :" + k)
	v, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(v)
	return string(v)
}
