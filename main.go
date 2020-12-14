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

type FfufCmdLine struct {
	CommandLine string `json:"commandline"`
}

type FfufTime struct {
	Time string `json:"time"`
}

//
type FfufConfig struct {
	Config FfufMetaData `json:"config"`
}

type FfufInputProviders struct {
	Keyword string `json:"keyword"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

// ffuf json 'commandline' struct
type FfufMetaData struct {
	URL        string `json:"url"`
	Method     string `json:"method"`
	Outputfile string `json:"outputfile"`
	//Time    string `json:"time"`
	InputProviders []FfufInputProviders `json:"inputproviders"`
}

// Array of ffuf {"results":..}
type FfufResults struct {
	Results []FfufResult `json:"results"`
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
func parseResults(byteVal []byte) {
	var commandline FfufCmdLine
	json.Unmarshal(byteVal, &commandline)

	var time FfufTime
	json.Unmarshal(byteVal, &time)

	var metadata FfufConfig
	json.Unmarshal(byteVal, &metadata)

	var results FfufResults
	json.Unmarshal(byteVal, &results)

	if debugmode {
		fmt.Printf("[debug] command line: %v\n", commandline.CommandLine)
		fmt.Printf("[debug] target: %v\n", metadata.Config.URL)
		fmt.Printf("[debug] method: %v\n", metadata.Config.Method)
		fmt.Printf("[debug] wordlist: %v\n", metadata.Config.InputProviders[0].Value)
		fmt.Printf("[debug] outputfile: %v\n", metadata.Config.Outputfile)
		fmt.Printf("[debug] time: %v\n", time.Time)
		fmt.Printf("[debug] url: %v\n", results.Results[0].URL)
	}

	fmt.Printf("%v %v %v %v %v\n", results.Results[0].URL, results.Results[0].Status, results.Results[0].Length, results.Results[0].Words, results.Results[0].Lines)

}

// set debug mode
func setDebugMode() {
	debug_str := os.Getenv("DEBUG")

	fmt.Printf("mode=%v\n", debug_str)
	if debug_str == "true" {
		fmt.Println("[*] Debug mode set to true.")
		debugmode = true
	} else {
		debugmode = false
	}
}
