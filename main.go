package main

import "fmt"
import "bufio"
import "os"
import "log"

//import "regexp"
import "strings"
import "path"

/*
 * Goal: Print a list of all known endpoints and the wordlist run against them
 * 1. take file as input
 * 2. parse the name of the file and collect the wordlist and the target
 *		gobuster-directory-list-2.3-small.txt-https.jss-dev.myseek.xyz.txt
 *		[ type ] [        wordlist          ] [          targ        ]
 *		... what to do about extensions?
 *			if str[] == 'gobuster-ext' {} else {}
 * 3. ...
 */
func main() {
	// hard code the input file for now
	input_file_path := "test-data/gobuster-directory-list-2.3-small.txt-https.getpwnd.com.txt"

	fmt.Println("input file: ", input_file_path)

	// get the file name from the full path
	filestr := path.Base(input_file_path)
	println("filestr: ", filestr)

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
	println("wordlist: ", wordlist)

	// open the file and print the contents, line by line
	file, err := os.Open(input_file_path)
	// unable to open the file for reading
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		urlstr := strings.Split(line, " ")
		println("  url: ", urlstr[0])
		//fmt.Println("  " + scanner.Text())
	}

	// no idea what this does
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
