package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Result struct {
	Target    string
	Wordlist  string
	Endpoints []Endpoint
	Tool      string
	Filepath  string
	Filename  string
}

type Endpoint struct {
	EP     string
	Status string
	Length string
}

// ffuf structs
type FfufResults struct {
	Results []FfufResult `json:"results"`
}

type FfufResult struct {
	Input    string `json:"input"`
	Position string `json:"position"`
	Status   int    `json:"status"`
	Length   int    `json:"length"`
	Words    string `json:"words"`
}

type FfufCmdLine struct {
	CommandLine string `json:"commandline"`
}

var results []Result

func main() {
	// Support for modes:
	//   1. endpoints (ep) - essentially cat; show all endpoints discovered
	//   2. targs - show all targeted endpoints and the wordlists used against them
	//   3. tree - show tree view of all discovered endpoints

	flag.Parse()

	mode := flag.Arg(0)

	// check we have legit mode
	if mode != "targs" && mode != "eps" && mode != "tree" {
		fmt.Fprintf(os.Stderr, "unknown mode %s\n", mode)
		return
	}

	// read from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fn_str := sc.Text()

		// only process existing files (not dirs)
		if fileExists(fn_str) {
			processFile(fn_str)
		}
	}

	switch mode {
	case "targs":
		ep_map := make(map[string][]string)

		// iterate through the results
		for _, a_res := range results {
			// append the new wordlist to the existing list of wordlists
			ep_map[a_res.Target] = targetBuildWordlistList(ep_map[a_res.Target], a_res.Wordlist)

		}

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
		for key, value := range ep_map {
			//w.Init(os.Stdout, 5, 0, 1, ' ', tabwriter.AlignRight)
			str := fmt.Sprintf("%s\t%s\t", key, targetPrettyPrintWordlists(value))
			fmt.Fprintln(w, str)
		}
		w.Flush()
	case "eps":
		// get and print a unique list of all endpoints
		endpoints := processEndpoints()

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
		//fmt.Fprintln(w, "Endpoint\tStatus\tSize\t")
		for _, ep := range endpoints {
			// TODO - print the expanded endpoint
			str := fmt.Sprintf("%s\t(Status: %s)\t(Size: %s)\t", ep.EP, ep.Status, ep.Length)
			fmt.Fprintln(w, str)
		}
		w.Flush()
	}
}

// nicely print the list of wordlists
func targetPrettyPrintWordlists(list []string) string {
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

// unique append
func targetBuildWordlistList(existing []string, addition string) []string {
	// pretty append wordlist to list of wordlists
	if !sliceContains(existing, addition) {
		existing = append(existing, addition)
	}

	return existing
}

// check if a slice contains an entry
func sliceContains(slice []string, needle string) bool {
	// Check if a slice contains a value

	for _, str := range slice {
		if str == needle {
			return true
		}
	}

	return false
}

// parse the wordlist from the filename
func resultGetWordlistName(raw string) string {
	// get the file name from the full path
	filestr := path.Base(raw)

	// Operation:
	//   1. ditch [0]
	//   2. collect [1] -> [N] , where [N] equals 'x.txt'

	// break down the string on the '-' character
	s := strings.Split(filestr, "-")

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

// parse the target from the filename
func resultGetTarget(filename string) string {
	// extract the target
	//   caveat: the target ended with a slash
	//     ie .com/dir/ becomes dir_.txt
	slice := strings.Split(filename, "-")

	// In the case of ffuf, the results only contain the position.
	// So if the target was targ.com/management, then the results
	// will only contain the discovered file/dir under /management.
	// As we need the whole thing for targs mode, we need to parse
	// the actualy target from the json's commandline value. It
	// would also be more accurate to pull the 'common denominator'
	// from gobuster, instead of relying on the file's name.

	return slice[len(slice)-1]
}

// parse the tool used from the filename
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

// parse the filename to extract the target, tool, and the wordlist
func processFile(filepath string) bool {
	var res Result

	res.Filepath = filepath
	res.Filename = path.Base(filepath)

	// ensure we're only going after output files
	if path.Ext(res.Filename) == ".txt" {
		// parse out the wordlist
		res.Wordlist = resultGetWordlistName(res.Filename)

		// parse out the tool
		res.Tool = resultGetTool(res.Filename)

		// parse out the target
		// TODO - remove trailing '_'
		res.Target = resultGetTarget(res.Filename)

		results = append(results, res)
	}

	return true
}

// check if file exists and is not directory
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

/* --- mode: eps --- */

// return a unique list of all endpoints
func processEndpoints() []Endpoint {
	all_eps := []string{}
	retEP := []Endpoint{}

	// parse each result's endpoints
	for i := 0; i < len(results); i++ {
		if results[i].Tool == "ffuf" {
			results[i].Endpoints = readInFfufEndpoints(results[i].Filepath)
		} else if results[i].Tool == "gobuster" {
			results[i].Endpoints = readInGobusterEndpoints(results[i].Filepath)
		}
	}

	// iterate through each result and add unseen endpoints
	for _, res := range results {
		for _, ep := range res.Endpoints {
			if !sliceContains(all_eps, ep.EP) {
				all_eps = append(all_eps, ep.EP)
				retEP = append(retEP, ep)
			}
		}
	}

	return retEP
}

// gather all endpoints from a result
func readInGobusterEndpoints(in_file string) []Endpoint {
	endpoints := []Endpoint{}

	// open the file and print the contents, line by line
	file, err := os.Open(in_file)

	// unable to open the file for reading
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read in the file line by line, append new eps to slice
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ep := Endpoint{}

		s := strings.Split(line, " ")
		ep.EP = epsParseRaw(s[0])
		ep.Status = strings.Replace(s[2], ")", "", 1)
		ep.Length = strings.Replace(s[4], "]", "", 1)

		endpoints = append(endpoints, ep)
	}

	// no idea what this does
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return endpoints
}

func epsParseRaw(ep_raw string) string {
	u, err := parseURL(ep_raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
	}

	return u.Path
}

func parseURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		return url.Parse("http://" + raw)
	}

	return u, nil
}

// read json blob file and parse
func readInFfufEndpoints(in_file string) []Endpoint {
	endpoints := []Endpoint{}

	jsonFile, err := os.Open(in_file)

	if err != nil {
		fmt.Println(err)
	}

	// do stuff
	defer jsonFile.Close()

	byteVal, _ := ioutil.ReadAll(jsonFile)

	// extract the results json from ffuf metadata
	var fResults FfufResults
	json.Unmarshal(byteVal, &fResults)

	// extract the full target from the ffuf commandline json metadata
	var fCmdLine FfufCmdLine
	json.Unmarshal(byteVal, &fCmdLine)
	target := ffufParseTarget(fCmdLine.CommandLine)

	var ep Endpoint
	for i := 0; i < len(fResults.Results); i++ {
		ep.EP = fResults.Results[i].Input
		// if a target was extracted, then manipulate the 'input's
		if target != "" {
			// concatenate the target + the input
			ep.EP = target + ep.EP
			// parse the path
			ep.EP = epsParseRaw(ep.EP)
		}

		ep.Status = strconv.Itoa(fResults.Results[i].Status)
		ep.Length = strconv.Itoa(fResults.Results[i].Length)
		//fmt.Println(target, fResults.Results[i].Input, " ", fResults.Results[i].Status)
		endpoints = append(endpoints, ep)
	}

	return endpoints
}

// parse the url/target from the 'commandline' json variable
func ffufParseTarget(cmd string) string {
	var t string

	args := strings.Split(cmd, " ")
	for _, arg := range args {
		if strings.Contains(arg, "FUZZ") {
			t = strings.Replace(arg, "FUZZ", "", 1)
		}
	}

	return t
}
