package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	kv "github.com/dcai/kv/src"
)

const JSON_DB_FILEPATH = "data.json"

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {
	log.SetPrefix("[go-kv] ")

	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatalln("Missing action. Try 'kv set key value' or 'kv get key'")
	}

	action := os.Args[1]

	kv.InitJsonDB(JSON_DB_FILEPATH)

	switch action {
	case "set":
		setCmd := flag.NewFlagSet("set", flag.ExitOnError)
		setCmd.Parse(os.Args[2:])

		var value string
		if isInputFromPipe() {
			bytes, err := io.ReadAll(os.Stdin)

			if err != nil {
				log.Fatal(err)
			}
			value = string(bytes)
		}
		key := setCmd.Args()[0]
		if len(setCmd.Args()) >= 2 {
			value = setCmd.Args()[1]
		}
		if value != "" {
			kv.SetValue(JSON_DB_FILEPATH, key, value)
		} else {
			log.Fatalln("Missing value. Try 'kv set key value'")
		}
	case "get":
		getCmd := flag.NewFlagSet("get", flag.ExitOnError)
		getAllRows := getCmd.Bool("all", false, "Get all keys and values")
		getCmd.Parse(os.Args[2:])
		var key string
		if len(getCmd.Args()) > 0 {
			key = getCmd.Args()[0]
		}
		if key != "" {
			value, err := kv.GetValue(JSON_DB_FILEPATH, key)
			if err == nil {
				fmt.Print(value)
			}
		} else {
			if *getAllRows {
				kv.PrintItemsRaw(JSON_DB_FILEPATH)
			} else {
				kv.PrintItemsInTable(JSON_DB_FILEPATH)
			}
		}
	case "rm":
		rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)
		rmCmd.Parse(os.Args[2:])
		var key string
		if len(rmCmd.Args()) > 0 {
			key = rmCmd.Args()[0]
		}
		kv.DeleteItem(JSON_DB_FILEPATH, key)
	case "rename":
		renameCmd := flag.NewFlagSet("rename", flag.ExitOnError)
		renameCmd.Parse(os.Args[2:])
		if len(renameCmd.Args()) >= 2 {
			oldkey := renameCmd.Args()[0]
			newkey := renameCmd.Args()[1]
			kv.RenameKey(JSON_DB_FILEPATH, oldkey, newkey)
		}

	default:
		os.Exit(1)
	}
}
