package main

import (
	"fmt"
	"time"
)

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
		//wurds := search5(t.wurds)
		reply := Thought{"wurds", 99}
		neurons <- reply
	}
}
