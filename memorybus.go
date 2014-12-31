package main

import (
	"fmt"

	_ "github.com/lib/pq"
	//"strings"
)

func (b *Bot) memoryBus(memory chan Event) {
	for m := range memory {
		//fmt.Printf("MEMORY:: when: %s // what: %s // ", m.When, m.What)
		b.CurrentStimulii = append(b.CurrentStimulii, m)
		fmt.Println("Memories LEN: ", len(b.CurrentStimulii))
	}
}

func (b *Bot) store(thing *Thing) error {
	count, err := b.Db.Update(thing)
	if err != nil {
		fmt.Println("Beep! error inserting to DB  :", err.Error())
		return err
	}
	fmt.Println("Updated ", count, " rows..")
	return nil
}

func (b *Bot) retrieveMemory(thing *Thing) (string, error) {
	err := b.Db.SelectOne(thing, "select * from thing where name=$1", thing.Name)
	if err != nil {
		fmt.Println("ERRRRX:", err.Error())
		fmt.Println(thing.Name, "doesnt exist - now creating..")
		err = b.Db.Insert(thing)
		if err != nil {
			fmt.Println("ERrrr inserting to DB:", err.Error())
			return "", err
		}
		return "didn't exist - now created", nil
	} else {
		return "woop! exists", nil
	}
}
