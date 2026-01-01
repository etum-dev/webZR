package main

import (
	"bufio"

	"fmt"
	"os"

	"github.com/etum-dev/WebZR/scan"
	"github.com/etum-dev/WebZR/utils"
)

// Optional file, args and/or stdin and scans domains concurrently
func processIn(inputFile string, args []string, hasStdin bool) {
	outputFile := "scan_results.json"
	// TODO: guard this collector with context cancellation + progress logging for long scans.

	if len(args) > 0 {
		for _, a := range args {
			domain := utils.CheckDomain(a)
			flags := utils.FlagParse()
		}
	}

	// Process file input
	if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file %s: %v\n", inputFile, err)
		} else {
			defer f.Close()
			sc := bufio.NewScanner(f)

			for sc.Scan() {
				//domain := utils.CheckDomain(sc.Text())
				fmt.Println(utils.CheckDomain(sc.Text()))
				// jobs <- domain
			}
			if err := sc.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", inputFile, err)
			}
		}
	}
	// Process piped stdin if present
	if hasStdin {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			domain := utils.CheckDomain(sc.Text())
			// add to wg
			fmt.Println(domain)
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		}
	}
	/*
		if totalResults > 0 {
			fmt.Printf("\n[*] Scan complete! Found %d results\n", totalResults)
		} else {

			fmt.Println("\no(TヘTo) くぅ No WebSocket connections found")
		}
	*/
}

func doScan(d string, flags string) /*[]utils.ScanResult*/ {
	// TODO: return richer telemetry (latency, error info)
	// TODO: finegrain flags to select type of scans
	cspEndpoints := scan.CSPSearch(d)

	// Test CSP endpoints and convert to ScanResults
	var cspResults []utils.ScanResult
	if len(cspEndpoints) > 0 {
		fmt.Printf("[CSP] Testing %d potential WebSocket endpoints for %s:\n", len(cspEndpoints), d)
		for _, endpoint := range cspEndpoints {
			fmt.Printf("  - Testing %s\n", endpoint)
			// Test if the endpoint actually responds to WebSocket requests
			result := scan.SendConnRequest(endpoint)
			if result != nil && result.Success {
				cspResults = append(cspResults, *result)
				fmt.Printf(" !! CSP endpoint confirmed: %s\n", endpoint)
			}
		}
	}
}

func main() {
	flags := utils.FlagParse()
	if flags != nil {
		fmt.Println("pizza")
	}

}
