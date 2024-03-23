package main

import (
	"github.com/jedib0t/go-pretty/v6/table"

	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"time"
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

// func writeJson(filename string, store KeyValueStore) {
// 	// data, err = json.Marshal(store)
// 	// if err != nil {
// 	//     log.Fatal(err)
// 	// }
// 	// err = ioutil.WriteFile(filename, data, 0644)
// 	// if err != nil {
// 	//     log.Fatal(err)
// 	// }
// }

func readFileAsString(filename string) string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
func createItem(key string, value string) *keyValueItem {
	now := time.Now().Format(time.RFC3339)
	base64Value := base64.StdEncoding.EncodeToString([]byte(value))
	checksum := stringChecksum(value)
	return &keyValueItem{
		Key:      key,
		Value:    base64Value,
		Created:  now,
		Updated:  now,
		Checksum: checksum,
	}
}

func updateItem(item *keyValueItem, value string) {
	now := time.Now().Format(time.RFC3339)
	base64Value := base64.StdEncoding.EncodeToString([]byte(value))
	checksum := stringChecksum(value)
	item.Value = base64Value
	item.Updated = now
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
			updateItem(&store.Items[index], value)
			break
		}
	}
	if !found {
		item := createItem(key, value)
		store.Items = append(store.Items, *item)
	}

	SaveStore(filename, &store)
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

func PrintAllItems(filename string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value"})
	str := readFileAsString(filename)
	store := keyValueStore{}
	json.Unmarshal([]byte(str), &store)
	for _, v := range store.Items {
		t.AppendRows([]table.Row{
			{v.Key, base64decode(v.Value)},
		})

	}
	t.AppendSeparator()

	t.Render()

}
