package main

import "fmt"
import "bufio"
import "os"
import "log"

//import "regexp"
import "strings"
import "path"

/*
 * Goal: Print a list of all known endpoints and the wordlist run against them
 * 1. Take a list of files from stdin
 * 2. parse the name of the file and collect the wordlist and the target
 *		... and the extension parameters?
 *
 *		gobuster-directory-list-2.3-small.txt-https.jss-dev.myseek.xyz.txt
 *		[ type ] [        wordlist          ] [          targ        ]
 *		... what to do about extensions?
 *			if str[] == 'gobuster-ext' {} else {}
 * 3. For all files, store a unique list of endpoints and the wordlist used
 * 4. Print a unique list of endpoints along with the each wordlist the endpoint was found in
 */

var endpoint_map = make(map[string][]string)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		/*
			u, err := getWordlistName(sc.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse failure %s\n", err)
				continue
			}
		*/
		wordlist := getWordlistName(sc.Text())
		//fmt.Print(wordlist, "\n")
		endpoints := readFile(sc.Text())
		//fmt.Print(endpoints)

		for _, str := range endpoints {
			addEntry(str, wordlist)
		}

	}
	for key, value := range endpoint_map {
		list := prettyList(value)
		fmt.Println(key, "(", list, ")")
	}

}

func getWordlistName(raw string) string {
	// get the file name from the full path
	filestr := path.Base(raw)
	//fmt.Println("filestr: ", filestr)

	/*
	 * ditch [0]
	 * collect [1] -> [N] , where [N] equals 'x.txt'
	 */

	// break down the string on the '-' character
	s := strings.Split(filestr, "-")
	//fmt.Println(s)

	// extract the wordlist (hacky)
	wordlist := ""
	get_wordlist := false
	// staring from 1 to skip the 'gobuster-' portion of the string
	for i := 1; i < len(s); i++ {
		//fmt.Println("slice: ", s[i])
		if get_wordlist == false {
			wordlist += s[i] + "-"
			if strings.Contains(s[i], ".txt") {
				//println("gotcha")
				get_wordlist = true
			}
		}
	}

	// remove the trailing '-'
	wordlist = strings.TrimSuffix(wordlist, "-")

	return wordlist
}

func readFile(input_file_path string) []string {
	endpoints := []string{}
	// open the file and print the contents, line by line
	file, err := os.Open(input_file_path)
	// unable to open the file for reading
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		urlstr := strings.Split(line, " ")
		//println("  url: ", urlstr[0])
		endpoints = append(endpoints, urlstr[0])
		//fmt.Println("  " + scanner.Text())
	}

	// no idea what this does
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return endpoints
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
