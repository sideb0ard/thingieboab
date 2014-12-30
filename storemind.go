package main

import (
	"fmt"

	_ "github.com/lib/pq"
	//"strings"
)

func (b *Bot) storeMind(memory chan Event) {
	for m := range memory {
		//fmt.Printf("MEMORY:: when: %s // what: %s // ", m.When, m.What)
		b.CurrentStimulii = append(b.CurrentStimulii, m)
		fmt.Println("Memories LEN: ", len(b.CurrentStimulii))
	}
}

//func (b *Bot) retrieveMind(thing Thing) {
//  if b.CurrentStimulii
//}
