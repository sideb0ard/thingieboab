package main

import (
	"fmt"
	"time"
)

type Thought struct {
	wurds string
	mood  int
}

func main() {
	fmt.Println("THINGS ARE GO!")
	neurons := make(chan Thought)
	done := make(chan bool)
	go uppermind(neurons)
	go lowermind(neurons)
	initThought := Thought{"HULLO", 100}
	neurons <- initThought
	_ = <-done
	fmt.Println("nite nite.")
}

func uppermind(neurons chan Thought) {
	for t := range neurons {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(t)
		reply := Thought{"UPPER MIND SHIZZLE IT", 99}
		neurons <- reply
	}
}

func lowermind(neurons chan Thought) {
	fmt.Println("STARTED!")
	for t := range neurons {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println(t)
		reply := Thought{"DEEEEP THOUGHT VIBEZZZ", 99}
		neurons <- reply
	}
}
