package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	// read filenames from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if debugmode {
			fmt.Printf("[debug] filename: %v\n", sc.Text())
		}
		fn_str := sc.Text()

		// read & process file
		if fileExists(fn_str) {
			byteVal := processFile(fn_str)
			parseResults(byteVal)
		}
	}
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
func parseResults(byteVal []byte) NavEndpoints {
	var endpoints NavEndpoints

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

	endpoints.Endpoint = output.Results[0].URL
	endpoints.Status = output.Results[0].Status
	endpoints.Length = output.Results[0].Length
	endpoints.Words = output.Results[0].Words
	endpoints.Lines = output.Results[0].Lines

	if debugmode {
		fmt.Printf("%v [Status: %v, Size: %v, Words: %v, Lines: %v]\n", endpoints.Endpoint, endpoints.Status, endpoints.Length, endpoints.Words, endpoints.Lines)
	}

	return endpoints
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
