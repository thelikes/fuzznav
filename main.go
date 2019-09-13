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
 * 1. Take a list of files from stdin
 * 2. parse the name of the file and collect the wordlist and the target
 *		... and the extension parameters?
 *
 *		gobuster-directory-list-2.3-small.txt-https.jss-dev.myseek.xyz.txt
 *		[ type ] [        wordlist          ] [          targ        ]
 *		... what to do about extensions?
 *			if str[] == 'gobuster-ext' {} else {}
 * 3. For all files, store a unique list of endpoints and the wordlist used
 * 4. Print a unique list of endpoints along with the each wordlist the endpoint was found in
 */
func main() {
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		/*
			u, err := parseFileName(sc.Text())
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse failure %s\n", err)
				continue
			}
		*/
		parseFileName(sc.Text())
		readFile(sc.Text())
	}
}

func parseFileName(raw string) {
	// get the file name from the full path
	filestr := path.Base(raw)
	fmt.Println("filestr: ", filestr)

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

}

func readFile(input_file_path string) {
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
