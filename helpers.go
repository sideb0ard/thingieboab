package main

import (
	"fmt"
	"log"
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

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
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
func isSubject(wurd string) bool {
	sjre, _ := regexp.Compile("([a-zA-Z0-9]).*SBJ")
	fmt.Println("%q\n", sjre.FindStringSubmatch(wurd))
	return true
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
