package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ffuf json struct
type FfufCommandLine struct {
	Commandline string
}

func main() {
	// read filenames from stdin
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		fmt.Printf("[debug] filename: %v\n", sc.Text())
		fn_str := sc.Text()

		if fileExists(fn_str) {
			processFile(fn_str)
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

// parse the filename to extract the target, tool, and the wordlist
func processFile(filepath string) {
	// read file
	jsonFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
	}
	byteVal, _ := ioutil.ReadAll(jsonFile)

	var cmdline FfufCommandLine
	json.Unmarshal(byteVal, &cmdline)
	fmt.Printf("[debug] command line: %v\n", cmdline)

	defer jsonFile.Close()
}
