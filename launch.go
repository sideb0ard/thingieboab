package main

import (
	"fmt"
)

func main() {
	//if len(os.Args) == 1 {
	//	log.Fatal("GIMME WURRRRRDZ")
	//}
	//fmt.Println("THINGS ARE GO!")
	//tokenz := tokenizer(strings.Join(os.Args[1:], " "))
	//for i := range tokenz {
	//	fmt.Println("TOKEN: ", tokenz[i])
	//}
	//	startwurds := os.Args[1]
	//	neurons := make(chan Thought)
	//	done := make(chan bool)
	//	go lowermind(neurons)
	//	go uppermind(neurons)
	//	initThought := Thought{startwurds, 100}
	//	neurons <- initThought
	//	_ = <-done
	var b Bot
	b.name = "AIGOR"
	b.mood = 100
	b.innit()
	b.talkPerson()
	fmt.Println("nite nite.")
}
