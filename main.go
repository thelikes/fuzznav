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
var debugmode bool

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

type NavEndpoints struct {
	Endpoint string
	Status   int
	Length   int
	Words    int
	Lines    int
}

func main() {
	// proces debug mode
	setDebugMode()

	// slice to store slices of NavEndpoints
	var results [][]NavEndpoints

	// read and parse each filename from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if debugmode {
			fmt.Printf("[debug] filename: %v\n", sc.Text())
		}
		fn_str := sc.Text()

		// read & process results from one file
		if fileExists(fn_str) {
			byteVal := processFile(fn_str)
			results = append(results, parseResults(byteVal))
			//fmt.Printf("len(results) = %v\n", len(results))
		}
	}

	endpointsMap(results)
}

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

// collect values
func parseResults(byteVal []byte) []NavEndpoints {
	var endpoints NavEndpoints
	var endpointResults []NavEndpoints

	var output FfufOutput
	json.Unmarshal(byteVal, &output)

	if debugmode {
		fmt.Printf("[debug] command line: %v\n", output.CommandLine)
		fmt.Printf("[debug] target: %v\n", output.Config.URL)
		fmt.Printf("[debug] method: %v\n", output.Config.Method)
		fmt.Printf("[debug] wordlist: %v\n", output.Config.InputProviders[0].Value)
		fmt.Printf("[debug] outputfile: %v\n", output.Config.Outputfile)
		fmt.Printf("[debug] time: %v\n", output.Time)
		fmt.Printf("[debug] url: %v\n", output.Results[0].URL)
	}

	//fmt.Printf("len: %v\n", len(output.Results))
	for i := 0; i < len(output.Results); i++ {
		endpoints.Endpoint = output.Results[i].URL
		endpoints.Status = output.Results[i].Status
		endpoints.Length = output.Results[i].Length
		endpoints.Words = output.Results[i].Words
		endpoints.Lines = output.Results[i].Lines

		if debugmode {
			fmt.Printf("%v [Status: %v, Size: %v, Words: %v, Lines: %v]\n", endpoints.Endpoint, endpoints.Status, endpoints.Length, endpoints.Words, endpoints.Lines)
		}

		endpointResults = append(endpointResults, endpoints)

	}

	return endpointResults
}

// print map of endpoints
func endpointsMap(results [][]NavEndpoints) {
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
func processEndpoints(results [][]NavEndpoints) []NavEndpoints {
	var allEndpoints []string
	var cleanResults []NavEndpoints

	// for each slice in results
	for _, res := range results {
		// for each result in results
		for _, ep := range res {
			// only store an endpoint if not already stored
			if !sliceContains(allEndpoints, ep.Endpoint) {
				allEndpoints = append(allEndpoints, ep.Endpoint)
				cleanResults = append(cleanResults, ep)
			}
		}
	}

	return cleanResults
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
		debugmode = true
	} else {
		debugmode = false
	}
}
