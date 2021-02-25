package util

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"text/tabwriter"
)

/*
 * Endpoint Processing Functions =====
 */

// print map of endpoints
func EndpointsMap(results [][]NavResults) {
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	endpoints := processEndpoints(results)

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
			// only process if endpoint is not null (no results)
			if len(ep.Endpoint) != 0 {
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
