package main

import "fmt"
import "bufio"
import "os"
import "log"

/*
 * Goal: Print a list of all known endpoints and the wordlist run against them
 * 1. take file as input
 * 2. parse the name of the file and collect the wordlist and the target
 * 3. ...
 */
func main() {
	fmt.Println("[*] Initializing...")

	// hard code the input file for now
	input_file_path := "test-data/run0/data.txt"

	fmt.Println(input_file_path)

	file, err := os.Open(input_file_path)
	// unable to open the file for reading
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println("  " + scanner.Text())
	}

	// no idea what this does
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
