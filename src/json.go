package kv

import (
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

func InspectStore(store *keyValueStore) {
	bytes, _ := json.MarshalIndent(store, "", "  ")
	log.Printf(string(bytes))
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

func GetStore(filename string) *keyValueStore {
	str := readFileAsString(filename)
	store := keyValueStore{}
	json.Unmarshal([]byte(str), &store)
	return &store
}

func DeleteItem(filename string, key string) {
	store := GetStore(filename)
	var found int
	for index, v := range store.Items {
		if v.Key == key {
			found = index
			break
		}
	}
	log.Printf("Deleting %s\n", key)

	if found >= 0 {
		store.Items = append(store.Items[:found], store.Items[found+1:]...)
		SaveStore(filename, store)
	}
}
func SetValue(filename string, key string, value string) {
	store := GetStore(filename)
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

	SaveStore(filename, store)
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
	store := GetStore(filename)
	for _, v := range store.Items {
		if v.Key == key {
			return base64decode(v.Value), nil
		}
	}
	return "", errors.New("key not found")
}

func RenameKey(filename string, oldKey string, newKey string) {
	store := GetStore(filename)
	found := false
	for index, v := range store.Items {
		if v.Key == oldKey {
			found = true
			store.Items[index].Key = newKey
		}
	}
	if found {
		// InspectStore(store)
		SaveStore(filename, store)
	}
}
