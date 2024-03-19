package main

import "fmt"
import "rsc.io/quote"
import "errors"

func Hello(name string) (string, error) {
	if name == "" {
        return "", errors.New("ERROR: empty name")

	}
	fmt.Println(quote.Go())

	return fmt.Sprintf("hello, %s", name), nil
}
