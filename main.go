package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const JSON_DB_FILEPATH = "data.json"

func main() {
	log.SetPrefix("[go-kv] ")
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	// updateIfExists := setCmd.Bool("u", false, "update if exists")
	// deleteKey := setCmd.Bool("D", false, "delete key")
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	// name := flag.String("name", "@dcai", "Name to greet")
	// ignore := flag.Bool("ignore", false, "should ignore")
	// printPersonInfo(name, ignore)

	flag.Parse()

	action := os.Args[1]

	InitJsonDB(JSON_DB_FILEPATH)

	switch action {
	case "set":
		setCmd.Parse(os.Args[2:])
		key := setCmd.Args()[0]
		value := setCmd.Args()[1]
		SetValue(JSON_DB_FILEPATH, key, value)
		fmt.Printf("key: %s, value: %s\n", key, value)
	case "get":
		getCmd.Parse(os.Args[2:])
		var key string
		if len(getCmd.Args()) > 0 {
			key = getCmd.Args()[0]
		}
		if key != "" {
			value, err := GetValue(JSON_DB_FILEPATH, key)
			if err == nil {
				fmt.Print(value)
			}
		} else {
			PrintAllItems(JSON_DB_FILEPATH)
		}
	default:
		fmt.Println("default")
		os.Exit(1)
	}

}
