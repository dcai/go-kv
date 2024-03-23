package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

const JSON_DB_FILEPATH = "data.json"

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {
	log.SetPrefix("[go-kv] ")
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAllRows := getCmd.Bool("all", false, "Get all keys and values")

	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatalln("Missing action. Try 'go-kv set key value' or 'go-kv get key'")
	}

	action := os.Args[1]

	InitJsonDB(JSON_DB_FILEPATH)

	switch action {
	case "set":
		setCmd.Parse(os.Args[2:])

		var value string
		if isInputFromPipe() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				value += scanner.Text() + "\n"
			}
		}
		key := setCmd.Args()[0]
		if len(setCmd.Args()) >= 2 {
			value = setCmd.Args()[1]
		}
		if value != "" {
			SetValue(JSON_DB_FILEPATH, key, value)
		} else {
			log.Fatalln("Missing value. Try 'go-kv set key value'")
		}
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
			if *getAllRows {
				PrintItemsRaw(JSON_DB_FILEPATH)
			} else {
				PrintItemsInTable(JSON_DB_FILEPATH)
			}
		}
	case "rm":
		fmt.Println("TODO - rm")
	case "rename":
		fmt.Println("TODO - rename")
	case "import":
		fmt.Println("TODO - import")
	default:
		fmt.Println("TODO - default")
		os.Exit(1)
	}
}
