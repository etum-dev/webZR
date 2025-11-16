package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/etum-dev/WebZR/basicfuzz"
	"github.com/etum-dev/WebZR/scan"
	"github.com/etum-dev/WebZR/utils"
)

// worker processes domains from the jobs channel and sends results to the results channel
func worker(id int, jobs <-chan string, results chan<- []utils.ScanResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for domain := range jobs {
		fmt.Printf("[Worker %d] Scanning domain: %s\n", id, domain)
		scanResults := doScan(domain)
		results <- scanResults
	}
	fmt.Printf("[Worker %d] Finished\n", id)
}

// processIn consumes optional file, args and/or stdin and scans domains concurrently
func processIn(inputFile string, args []string, hasStdin bool, numWorkers int) {
	// Channels for concurrent processing
	jobs := make(chan string, 100)
	results := make(chan []utils.ScanResult, 100)
	var wg sync.WaitGroup

	// Start worker pool
	fmt.Printf("[*] Starting %d workers...\n", numWorkers)
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	// Goroutine to collect all results
	allResults := []utils.ScanResult{}
	done := make(chan bool)

	go func() {
		for scanResults := range results {
			allResults = append(allResults, scanResults...)
		}
		done <- true
	}()

	// Send jobs to workers
	// Process positional args first (if any)
	if len(args) > 0 {
		for _, a := range args {
			domain := utils.CheckDomain(a)
			jobs <- domain
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
				domain := utils.CheckDomain(sc.Text())
				jobs <- domain
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
			jobs <- domain
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		}
	}

	close(jobs)
	wg.Wait()
	close(results)
	<-done

	// Write all results to JSON file
	if len(allResults) > 0 {
		outputFile := "scan_results.json"
		fmt.Printf("\n[*] Scan complete! Found %d results\n", len(allResults))
		if err := utils.WriteMultipleResults(outputFile, allResults); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write results: %v\n", err)
		}
	} else {
		fmt.Println("\no(TヘTo) くぅ No WebSocket connections found")
	}
}
func doScan(d string) []utils.ScanResult {
	cspEndpoints := scan.CSPSearch(d)
	if len(cspEndpoints) > 0 {
		fmt.Printf("[CSP] Potential WebSocket endpoints for %s:\n", d)
		for _, endpoint := range cspEndpoints {
			fmt.Printf("  - %s\n", endpoint)
		}
	}

	endpoints := scan.ScanEndpoint(d)
	subdomains := scan.ScanSubdomain(d)

	// collect all results
	results := []utils.ScanResult{endpoints}
	results = append(results, subdomains...)

	return results
}

func main() {
	// flags
	isTest := flag.Bool("test", false, "Use the local test server")
	fuzzOpt := flag.Int("fuzz", 0, "Fuzzing option. 1 for basic, 2 for custom header, 3 for mutation")
	fileOpt := flag.String("file", "", "Path to input file (optional)")
	singleDomain := flag.String("domain", "", "Single domain")
	workers := flag.Int("workers", 5, "Number of concurrent workers (default: 5)")

	flag.Parse()

	// detect positional args
	args := flag.Args()

	// detect whether there's piped stdin
	hasStdin := false
	if fi, err := os.Stdin.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			hasStdin = true
		}
	}

	if *fileOpt != "" || len(args) > 0 || hasStdin {
		processIn(*fileOpt, args, hasStdin, *workers)
	} else if *singleDomain != "" {
		d := utils.CheckDomain(*singleDomain)
		result := scan.SendConnRequest(d)

		if result != nil && result.Success {
			fmt.Printf("\nヾ(o･∀･)ﾉ ﾀﾞｰ!![SUCCESS] WebSocket connection details:\n")
			fmt.Printf("  URL: %s\n", result.URL)
			fmt.Printf("  Scheme: %s\n", result.Scheme)
			fmt.Printf("  Insecure: %v\n", result.Insecure)
			fmt.Printf("  Status: %d\n", result.StatusCode)
		} else {
			fmt.Printf("\n[-] No WebSocket connection available for %s\n", d)
		}
	} else {

		fmt.Println("Usage: provide -file=domains.txt, -domain=example.com, or domains as arguments")
	}

	if *isTest {
		basicfuzz.ServeWs()
	}

	switch *fuzzOpt {
	case 1:
		fmt.Println("Doing a simple fuzz.")
		basicfuzz.SimpleFuzz()
	case 0:
		fmt.Println("Not fuzzing")
	}

}
