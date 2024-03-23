package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type keyValueItem struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
	Checksum string `json:"checksum"`
}

type keyValueStore struct {
	Version string         `json:"version"`
	Items   []keyValueItem `json:"items"`
}

func stringChecksum(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InitJsonDB(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		jsonDbFile, err := os.Create(filename)
		item := createItem("testkey", "123")

		store := &keyValueStore{
			Version: "1.0",
			Items:   []keyValueItem{*item},
		}
		bytes, _ := json.MarshalIndent(store, "", "  ")
		jsonDbFile.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		defer jsonDbFile.Close()
	}
}

func SaveStore(filename string, store *keyValueStore) {
	bytes, _ := json.MarshalIndent(store, "", "  ")
	ioutil.WriteFile(filename, bytes, 0644)
}

func readFileAsString(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func now() string {
	return time.Now().Format(time.RFC3339)
}
func createItem(key string, value string) *keyValueItem {
	nowstr := now()
	base64Value := base64encode(value)
	checksum := stringChecksum(value)
	return &keyValueItem{
		Key:      key,
		Value:    base64Value,
		Created:  nowstr,
		Updated:  nowstr,
		Checksum: checksum,
	}
}

func updateItem(item *keyValueItem, value string) {
	nowstr := now()
	base64Value := base64encode(value)
	checksum := stringChecksum(value)
	item.Value = base64Value
	item.Updated = nowstr
	item.Checksum = checksum
}

func SetValue(filename string, key string, value string) {
	str := readFileAsString(filename)
	store := keyValueStore{}
	json.Unmarshal([]byte(str), &store)

	found := false
	for index, v := range store.Items {
		if v.Key == key {
			found = true
			log.Println("Updating key: " + key)
			updateItem(&store.Items[index], value)
			break
		}
	}
	if !found {
		log.Println("Adding key: " + key)
		item := createItem(key, value)
		store.Items = append(store.Items, *item)
	}

	SaveStore(filename, &store)
}

func base64encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64decode(str string) string {
	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func GetValue(filename string, key string) (string, error) {
	str := readFileAsString(filename)
	store := keyValueStore{}
	json.Unmarshal([]byte(str), &store)
	for _, v := range store.Items {
		if v.Key == key {
			return base64decode(v.Value), nil
		}
	}
	return "", errors.New("key not found")
}

func PrintItemsInTable(filename string) {
	t := table.NewWriter()
	t.SetAllowedRowLength(80)
	// t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value"})
	str := readFileAsString(filename)
	store := keyValueStore{}
	json.Unmarshal([]byte(str), &store)
	for _, v := range store.Items {
		value := strings.ReplaceAll(base64decode(v.Value), "\n", "")
		t.AppendRows([]table.Row{
			{v.Key, value},
		})

	}
	t.AppendSeparator()
	t.Render()
}

func PrintRow(item *keyValueItem) {
	value := base64decode(item.Value)
	sep := text.FgHiBlack.Sprint(":")
	fmt.Printf("KEY     %s %s\n", sep, text.FgHiRed.Sprint(item.Key))
	fmt.Printf("CREATED %s %s\n", sep, text.FgYellow.Sprint(item.Created))
	fmt.Printf("UPDATED %s %s\n", sep, text.FgYellow.Sprint(item.Updated))
	fmt.Printf("VALUE   %s %s\n", sep, text.FgGreen.Sprint(value))
}

func PrintItemsRaw(filename string) {
	str := readFileAsString(filename)
	store := keyValueStore{}
	json.Unmarshal([]byte(str), &store)

	for _, v := range store.Items {
		PrintRow(&v)
		fmt.Println(text.FgHiBlack.Sprint(strings.Repeat("=", 80)))
	}
}
