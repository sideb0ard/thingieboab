package main

import (
	"fmt"
	//	"time"
)

func main() {
	//bored_chan := make(chan bool)
	//done_chan := make(chan bool)
	//listen_chan := make(chan string)
	//mood_chan := make(chan int)
	//neurons_chan := make(chan Thought)

	var b Bot
	b.Name = "AIGOR"
	b.Mood = 100
	b.innit()

	//	go b.listen(listen_chan)
	//	go b.lowermind(mood_chan, neurons_chan)
	//	go b.moody(mood_chan)
	//	go b.uppermind(mood_chan, neurons_chan)
	//	go b.talkPerson(bored_chan, listen_chan, mood_chan, neurons_chan)
	//
	//	for {
	//		time.Sleep(60 * time.Second)
	//		b.think(bored_chan, mood_chan, neurons_chan)
	//	}

	fmt.Println("nite nite.")
}
