package main

import (
	"math/rand"
	"strings"

	_ "github.com/lib/pq"
	//"strings"
	"time"
)

func (b *Bot) yinMind(neurons chan Thought, memory chan Event) {
	for t := range neurons {
		time.Sleep(time.Duration(random(0, 1000)) * time.Millisecond)
		//fmt.Println("yin", t.Wurds)

		wurds := strings.Split(t.Wurds, " ")
		//var replywurds string
		for i := range wurds {
			j := rand.Intn(i + 1)
			wurds[i], wurds[j] = wurds[j], wurds[i]
		}

		mem := Event{time.Now(), "neuron received", &b.Thing}
		memory <- mem

		reply := Thought{strings.Join(wurds, " "), 99}
		//fmt.Println("YIN SENDING REPLY", reply)
		neurons <- reply
	}
}
