package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

type Prefix []string

func (p Prefix) String() string {
	return strings.Join(p, " ")
}

func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

type Chain struct {
	chain     map[string][]string
	prefixLen int
}

func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		//fmt.Println("KEY IS", key)
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

//func (c *Chain) Build(s string) {
//	fmt.Println("BUILD - gots a", s)
//	p := make(Prefix, c.prefixLen)
//	key := p.String()
//	c.chain[key] = append(c.chain[key], s)
//	p.Shift(s)
//}

func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			continue
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	fmt.Println("Returning WURDS", words)
	return strings.Join(words, " ")
}
