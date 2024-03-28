package kv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dcai/kv/src/colors"
)

func PrintRow(item *keyValueItem) {
	value := base64decode(item.Value)
	sep := colors.FgGray(":")
	fmt.Printf("KEY     %s %s\n", sep, colors.FgRed(item.Key))
	fmt.Printf("CREATED %s %s\n", sep, colors.FgYellow(item.Created))
	fmt.Printf("UPDATED %s %s\n", sep, colors.FgYellow(item.Updated))
	fmt.Printf("VALUE   %s %s\n", sep, colors.FgGreen(value))
}

func PrintItemsRaw(filename string) {
	store := GetStore(filename)

	for _, v := range store.Items {
		PrintRow(&v)
		fmt.Println(colors.FgGray(strings.Repeat("=", 80)))
	}
}

func PrintItemsInTable(filename string) {
	store := GetStore(filename)
	longest_col := 30
	for _, v := range store.Items {
		if len(v.Key) > longest_col {
			longest_col = len(v.Key)
		}
	}

	formatter := "%-" + strconv.Itoa(longest_col) + "s%-23s%-23s\n"
	fmt.Printf(colors.FgYellow(formatter), "Key", "Created", "Updated")
	for _, v := range store.Items {
		fmt.Printf(formatter, v.Key, v.Created, v.Updated)
		// value := strings.ReplaceAll(base64decode(v.Value), "\n", "")
	}
}
