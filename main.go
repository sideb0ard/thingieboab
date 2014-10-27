package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	//"os/exec"
	//"time"
)

const (
	CONN_HOST = ""
	CONN_PORT = "7474"
	CONN_TYPE = "tcp"
)

func main() {
	//bored_chan := make(chan bool)
	//done_chan := make(chan bool)
	//listen_chan := make(chan string)
	//mood_chan := make(chan int)
	//neurons_chan := make(chan Thought)
	//cmd := exec.Command("clear")
	//cmd.Stdout = os.Stdout
	//cmd.Run()

	var debug = flag.Bool("d", false, "debug - whether to print copious what-i-m-doing messages")
	flag.Parse()

	var keywurds Keywurds

	var b Bot
	b.Name = "AIGOR"
	b.Mood = 100
	b.Debug = *debug
	b.innit(&keywurds)

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		panic("Error Listening:" + err.Error())
	}
	conns := clientConns(l)
	for {
		go handleRequest(<-conns)
	}
}

func clientConns(l net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		fmt.Println("YAr, Chatbot Listening on " + CONN_HOST + ":" + CONN_PORT)
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Err accepting: ", err.Error())
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, conn.LocalAddr(), conn.RemoteAddr())
			ch <- conn
		}
	}()
	return ch
}

//go b.listen(listen_chan)
//go b.lowermind(mood_chan, neurons_chan)
//go b.moody(mood_chan)
//go b.uppermind(mood_chan, neurons_chan)
//go b.talkPerson(bored_chan, listen_chan, mood_chan, neurons_chan, keywurds)

//for {
//	time.Sleep(60 * time.Second)
//	//b.think(bored_chan, mood_chan, neurons_chan)
//}

func handleRequest(conn net.Conn) {
	b := bufio.NewReader(conn)

	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			fmt.Println("Errz reading: ", err.Error())
			break
		}
		//conn.Write([]byte(reply))
		conn.Write([]byte(line))
	}
}
