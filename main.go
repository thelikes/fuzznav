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

	/*
		for _, a_res := range results {
			resultPrettyPrint(a_res)
		}
	*/

	switch mode {
	case "targs":
		ep_map := make(map[string][]string)

		// iterate through the results
		for _, a_res := range results {
			// append the new wordlist to the existing list of wordlists
			ep_map[a_res.Target] = targetBuildWordlistList(ep_map[a_res.Target], a_res.Wordlist)

		}

		for key, value := range ep_map {
			fmt.Printf("%s:%s\n", key, targetPrettyPrintWordlists(value))
		}

	}
	// identify - identify tool used to generate results

	// parse - parse the file name and/or the file contents
	// print - print results based on mode
	//  either iterate through a list of targets and print
	//  or iterate though targets list and print endpoints
}

func targetPrettyPrintWordlists(list []string) string {
	// nicely print the list of wordlists
	ret := ""

	for i := 0; i < len(list); i++ {
		if i == len(list)-1 {
			ret += list[i]
		} else {
			ret = ret + list[i] + ","
		}
	}

	return ret
}

func targetBuildWordlistList(existing []string, addition string) []string {
	// pretty append wordlist to list of wordlists
	if !sliceContains(existing, addition) {
		existing = append(existing, addition)
	}

	//fmt.Println("existing:", existing)
	//fmt.Println("a_res.Wordlist:", addition)

	return existing
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

func resultGetWordlistName(raw string) string {
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
			// wordlist has to end in '.txt'
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

func resultGetTarget(filename string) string {
	// extract the target
	slice := strings.Split(filename, "-")

	return slice[len(slice)-1]

}

func resultGetTool(filename string) string {
	// this one is easy because the tool is the first index of delimeter '-'
	// gobuster-directory-list-2.3-small.txt-http___towers.att.com.txt
	str := strings.Split(filename, "-")
	var tool string

	switch str[0] {
	case "ffuf":
		tool = "ffuf"
	case "gobuster":
		tool = "gobuster"
	default:
		tool = "unkown"
	}

	// extract the target
	return tool

}

// parse the filename to extract the target and the wordlist
func processFile(filepath string) bool {
	var res Result

	res.Filename = path.Base(filepath)

	// ensure we're only going after output files
	if path.Ext(res.Filename) == ".txt" {
		//fmt.Println("Processing file", res.Filename)

		// parse out the wordlist
		res.Wordlist = resultGetWordlistName(res.Filename)

		// parse out the tool
		res.Tool = resultGetTool(res.Filename)

		// parse out the target
		res.Target = resultGetTarget(res.Filename)

		//fmt.Println("  Wordlist:", res.Wordlist)

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

func resultPrettyPrint(res Result) {
	fmt.Println("Filename:", res.Filename)
	fmt.Println("Wordlist:", res.Wordlist)
	fmt.Println("Target:", res.Target)
	fmt.Println("Tool:", res.Tool)
}
