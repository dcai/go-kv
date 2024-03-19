package main

import "fmt"
import "log"

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	message, err := Hello("@dcai")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(message)
}
