package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var conceptnet5 string = "http://conceptnet5.media.mit.edu/data/5.2/search"

func search5(phrase string) string {
	url := conceptnet5 + "?text=" + phrase
	fmt.Println("URL5 is " + url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
	return string(body)
}
