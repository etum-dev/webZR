package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/etum-dev/WebZR/scan"
	"github.com/etum-dev/WebZR/utils"
)

const outputFile = "scan_results.json"

// scanJob is the unit of work for each worker.
func scanJob(job utils.Job) utils.JobResult {
	domain := utils.CheckDomain(job.Domain)
	if domain == "" {
		return utils.JobResult{Job: job, Err: fmt.Errorf("empty domain")}
	}

	var results []utils.ScanResult

	// Run scans based on mode/flags
	mode := job.Flags.Mode
	if mode == "" {
		mode = "basic" // default
	}

	// Always try CSP first (fast, no brute force)
	if cspResults := scan.ScanCSP(domain); len(cspResults) > 0 {
		results = append(results, cspResults...)
	}

	// Endpoint scanning (basic mode)
	if mode == "basic" || mode == "aggressive" {
		if epResults := scan.ScanEndpoint(domain); len(epResults) > 0 {
			results = append(results, epResults...)
		}
	}

	// Subdomain scanning (aggressive mode only)
	if mode == "aggressive" {
		if subResults := scan.ScanSubdomain(domain); len(subResults) > 0 {
			results = append(results, subResults...)
		}
	}

	return utils.JobResult{
		Job:     job,
		Results: results,
	}
}

// pushes domains from flags, args, and stdin into workers its just like palworld fr
func enqueueInputs(flags *utils.Flags, extraArgs []string, pool *utils.WorkerPool) error {
	var firstErr error

	enqueue := func(domain string) {
		pool.Submit(utils.Job{Domain: domain, Flags: flags})
	}

	if flags.DomainInput != "" {
		if utils.IsFile(flags.DomainInput) {
			file, err := os.Open(flags.DomainInput)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				enqueue(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				firstErr = err
			}
		} else {
			enqueue(flags.DomainInput)
		}
	}

	for _, arg := range extraArgs {
		enqueue(arg)
	}

	if stdinHasData() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			enqueue(scanner.Text())
		}
		if err := scanner.Err(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}

func stdinHasData() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}

func main() {
	flags := utils.FlagParse()
	scan.SetVerbose(flags.Verbose)

	workerCount := runtime.NumCPU()
	if workerCount < 2 {
		workerCount = 2
	}

	pool := utils.NewWorkerPool(workerCount, workerCount*4, scanJob)

	resultBatches := make(chan []utils.ScanResult, workerCount)
	writerDone := utils.StreamResults(outputFile, resultBatches)

	go func() {
		for res := range pool.Results() {
			if res.Err != nil {
				fmt.Fprintf(os.Stderr, "scan error for %s: %v\n", res.Job.Domain, res.Err)
				continue
			}
			if len(res.Results) > 0 {
				resultBatches <- res.Results
			}
		}
		close(resultBatches)
	}()

	if err := enqueueInputs(flags, flag.Args(), pool); err != nil {
		fmt.Fprintf(os.Stderr, "input warning: %v\n", err)
	}

	pool.Close()
	total := <-writerDone
	fmt.Printf("[*] Scan complete. %d result(s) written to %s\n", total, outputFile)
}
