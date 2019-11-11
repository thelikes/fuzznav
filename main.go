package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type Result struct {
	Target   string
	Wordlist string
	Endpoint string
	Tool     string
	Filename string
}

var results []Result

func main() {
	/* Support for modes
	 * 1. endpoints (ep) - essentially cat; show all endpoints discovered
	 * 2. targs - show all targeted endpoints and the wordlists used against them
	 * 3. tree - show tree view of all discovered endpoints
	 */

	flag.Parse()

	mode := flag.Arg(0)

	// check we have legit mode
	if mode != "targs" && mode != "ep" && mode != "tree" {
		fmt.Fprintf(os.Stderr, "unknown mode %s\n", mode)
		return
	}

	// read from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fn_str := sc.Text()

		// only process existing files (not dirs)
		if fileExists(fn_str) {
			//fmt.Println("Found file", fn_str)
			processFile(fn_str)
		}
	}

	fmt.Println(results)

	switch mode {
	case "targs":
		fmt.Println("Printing targs")
	}
	// identify - identify tool used to generate results

	// parse - parse the file name and/or the file contents
	// print - print results based on mode
	//  either iterate through a list of targets and print
	//  or iterate though targets list and print endpoints
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

// parse the filename to extract the target and the wordlist
func processFile(filepath string) bool {
	var res Result

	res.Filename = path.Base(filepath)

	// ensure we're only going after output files
	if path.Ext(res.Filename) == ".txt" {
		fmt.Println("Processing file", res.Filename)

		// parse out the wordlist
		res.Wordlist = getWordlistName(res.Filename)

		fmt.Println("  Wordlist:", res.Wordlist)

		results = append(results, res)

	}

	return true
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
