package main

import "fmt"
import "strings"

var endpoint_map = make(map[string][]string)

func main() {

	// for endpoints in file, add wordlist
	addEntry("/index.php", "big.txt")
	addEntry("/index.php", "common.txt")
	addEntry("/index.php", "common.txt")
	addEntry("/index.php", "big.txt")
	addEntry("/upload", "big.txt")
	addEntry("/upload/files", "directories.txt")

	for key, value := range endpoint_map {
		list := prettyList(value)
		fmt.Println(key, "(", list, ")")
	}
}

func addEntry(endpoint, wordlist string) {
	/*
	 * Append new wordlist to the entry
	 */

	slice := endpoint_map[endpoint]
	if !sliceContains(slice, wordlist) {
		slice = append(slice, wordlist)
	}

	endpoint_map[endpoint] = slice
}

func sliceContains(slice []string, needle string) bool {
	/*
	 * Check if a slice contains a value
	 */

	for _, str := range slice {
		if str == needle {
			return true
		}
	}

	return false
}

func prettyList(list []string) string {
	/* Accept a slice and return a pretty string */
	var listlist strings.Builder

	for i := 0; i < len(list); i++ {
		listlist.WriteString(list[i])
		if i < len(list)-1 {
			listlist.WriteString(", ")
		}
	}

	return listlist.String()
}
