package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	name := flag.String("name", "@dcai", "Name to greet")
	ignore := flag.Bool("ignore", false, "should ignore")

	flag.Parse()

	url := fmt.Sprintf("https://httpbin.org/get?name=%s", *name)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	sb := string(body)
	log.Println(sb)

	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	message, err := Hello(*name)
	if err != nil {
		log.Fatal(err)
	}
	if *ignore {
		fmt.Println("I am ignoring you")
	} else {
		fmt.Println(message)
	}
}
