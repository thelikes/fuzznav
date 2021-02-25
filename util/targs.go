package util

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
)

/*
 * Target Processing Functions =====
 */

// print unique list of target alongside unique list of wordlists
func TargetsMap(results [][]NavResults) {
	targetMap := make(map[string][]string)

	for _, res := range results {
		for _, aRes := range res {
			targetMap[aRes.URL] = targetBuildWordlistList(targetMap[aRes.URL], filepath.Base(aRes.Wordlist))
		}
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	for key, value := range targetMap {
		str := fmt.Sprintf("%s\t%s\t", key, targetPrettyPrintWordlists(value))
		fmt.Fprintln(w, str)
	}
	w.Flush()
}

// only add the wordlist if not already present
func targetBuildWordlistList(existing []string, addition string) []string {
	if !sliceContains(existing, addition) {
		existing = append(existing, addition)
	}

	return existing
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
