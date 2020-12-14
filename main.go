package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

// global debug mode used with env variable
var DEBUGMODE bool

/*
 * Ffuf JSON Structs
 */

type FfufOutput struct {
	CommandLine string       `json:"commandline"`
	Time        string       `json:"time"`
	Config      FfufMetaData `json:"config"`
	Results     []FfufResult `json:"results"`
}

type FfufInputProviders struct {
	Keyword string `json:"keyword"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

// ffuf json 'commandline' struct
type FfufMetaData struct {
	URL            string               `json:"url"`
	Method         string               `json:"method"`
	Outputfile     string               `json:"outputfile"`
	InputProviders []FfufInputProviders `json:"inputproviders"`
}

// ffuf json 'results' struct - status, size, words, lines
type FfufResult struct {
	Input    string `json:"input"`
	Position int    `json:"position"`
	Status   int    `json:"status"`
	Length   int    `json:"length"`
	Words    int    `json:"words"`
	Lines    int    `json:"lines"`
	Host     string `json:"host"`
	URL      string `json:"url"`
}

/*
 * FuzzNav Structs
 */

type NavResults struct {
	Endpoint   string
	Status     int
	Length     int
	Words      int
	Lines      int
	URL        string
	Wordlist   string
	Outputfile string
}

/*
 * Main
 */

func main() {
	// proces debug mode
	setDebugMode()

	// slice to store slices of NavResults
	var results [][]NavResults

	// read and parse each filename from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if DEBUGMODE {
			fmt.Printf("[debug] filename: %v\n", sc.Text())
		}
		fn_str := sc.Text()

		// read & process results from one file
		if fileExists(fn_str) {
			byteVal := processFile(fn_str)
			results = append(results, parseResults(byteVal))
		}
	}

	endpointsMap(results)
	//targetsMap(results)
}

/*
 * Accept JSON byte stream as input, parse and extract desired values.
 * Returns slice of one or more Ffuf result objects and additional metadata.
 *
 * Note: Lazy sotring metadata for each Ffuf result (ie discovered endpoint), though may have
 *   some benefits. Good place to start for performance optimization.
 */

func parseResults(byteVal []byte) []NavResults {
	var navResults NavResults
	var navResultsSlice []NavResults

	var output FfufOutput
	json.Unmarshal(byteVal, &output)

	if DEBUGMODE {
		fmt.Printf("[debug] command line: %v\n", output.CommandLine)
		fmt.Printf("[debug] target: %v\n", output.Config.URL)
		fmt.Printf("[debug] method: %v\n", output.Config.Method)
		fmt.Printf("[debug] wordlist: %v\n", output.Config.InputProviders[0].Value)
		fmt.Printf("[debug] outputfile: %v\n", output.Config.Outputfile)
		fmt.Printf("[debug] time: %v\n", output.Time)
		fmt.Printf("[debug] url: %v\n", output.Results[0].URL)
	}

	navResults.URL = output.Config.URL
	navResults.Outputfile = output.Config.Outputfile
	navResults.Wordlist = output.Config.InputProviders[0].Value

	if len(output.Results) == 0 {
		// handle storage of metadata for output file with no results
		navResultsSlice = append(navResultsSlice, navResults)
	} else {
		for i := 0; i < len(output.Results); i++ {
			navResults.Endpoint = output.Results[i].URL
			navResults.Status = output.Results[i].Status
			navResults.Length = output.Results[i].Length
			navResults.Words = output.Results[i].Words
			navResults.Lines = output.Results[i].Lines

			if DEBUGMODE {
				fmt.Printf("%v [Status: %v, Size: %v, Words: %v, Lines: %v]\n", navResults.Endpoint, navResults.Status, navResults.Length, navResults.Words, navResults.Lines)
			}

			navResultsSlice = append(navResultsSlice, navResults)
		}
	}

	return navResultsSlice
}

/*
 * Endpoint Processing Functions =====
 */

// print map of endpoints
func endpointsMap(results [][]NavResults) {
	endpoints := processEndpoints(results)
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	var str string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	for _, ep := range endpoints {
		switch {
		case ep.Status >= 200 && ep.Status <= 299:
			str = fmt.Sprintf("%v\t[Status: %v, Size: %v, Words: %v, Lines: %v]\t", ep.Endpoint, green(ep.Status), ep.Length, ep.Words, ep.Lines)
		case ep.Status >= 300 && ep.Status <= 399:
			str = fmt.Sprintf("%v\t[Status: %v, Size: %v, Words: %v, Lines: %v]\t", ep.Endpoint, blue(ep.Status), ep.Length, ep.Words, ep.Lines)
		case ep.Status >= 400 && ep.Status <= 499:
			str = fmt.Sprintf("%v\t[Status: %v, Size: %v, Words: %v, Lines: %v]\t", ep.Endpoint, yellow(ep.Status), ep.Length, ep.Words, ep.Lines)
		case ep.Status >= 500 && ep.Status <= 599:
			str = fmt.Sprintf("%v\t[Status: %v, Size: %v, Words: %v, Lines: %v]\t", ep.Endpoint, red(ep.Status), ep.Length, ep.Words, ep.Lines)
		default:
			str = fmt.Sprintf("%v\t[Status: %v, Size: %v, Words: %v, Lines: %v]\t", ep.Endpoint, ep.Status, ep.Length, ep.Words, ep.Lines)
		}

		fmt.Fprintln(w, str)
	}

	w.Flush()
}

// return a unique list of all endpoints
func processEndpoints(results [][]NavResults) []NavResults {
	var allEndpoints []string
	var cleanResults []NavResults

	// for each slice in results
	for _, res := range results {
		// for each result in results
		for _, ep := range res {
			//fmt.Printf("endpoint (%v): %v\n", ep.Outputfile, ep.Endpoint)
			// only process if endpoint is not null (no results)
			if len(ep.Endpoint) != 0 {
				//fmt.Printf("No endpoint found (%v)\n", ep.Outputfile)
				// only store an endpoint if not already stored
				if !sliceContains(allEndpoints, ep.Endpoint) {
					allEndpoints = append(allEndpoints, ep.Endpoint)
					cleanResults = append(cleanResults, ep)
				}
			}
		}
	}

	return cleanResults
}

/*
 * Target Processing Functions =====
 */

func targetsMap(results [][]NavResults) {
	targetResults := processTargets(results)

	var str string

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	for _, target := range targetResults {
		str = fmt.Sprintf("%v\t", target.URL)
		fmt.Fprintln(w, str)
	}
	w.Flush()
}

func processTargets(results [][]NavResults) []NavResults {
	var allTargets []string
	var cleanResults []NavResults

	// for each slice in results
	for _, res := range results {
		// for each result in results
		for _, aRes := range res {
			// only store an endpoint if not already stored
			if !sliceContains(allTargets, aRes.URL) {
				allTargets = append(allTargets, aRes.URL)
				cleanResults = append(cleanResults, aRes)
			}
		}
	}

	return cleanResults
}

/*
 * Utility Functions =====
 */

// check if file exists and is not directory
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// ingest json to byte array
func processFile(filepath string) []byte {
	// read file
	jsonFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
	}
	byteVal, _ := ioutil.ReadAll(jsonFile)

	defer jsonFile.Close()

	return byteVal
}

// check if a slice contains an entry
func sliceContains(slice []string, needle string) bool {
	for _, str := range slice {
		if str == needle {
			return true
		}
	}

	return false
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

// set debug mode
func setDebugMode() {
	debug_str := os.Getenv("DEBUG")

	if debug_str == "true" {
		fmt.Println("[*] Debug mode set to true.")
		DEBUGMODE = true
	} else {
		DEBUGMODE = false
	}
}
