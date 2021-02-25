package util

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func ReadStdinAndParse() [][]NavResults {
	// proces debug mode
	setDebugMode()

	// command line args
	flag.Parse()

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

	return results
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
