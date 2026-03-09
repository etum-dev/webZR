package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	config "github.com/etum-dev/WebZR/internal/flaginput"
	scanjob "github.com/etum-dev/WebZR/internal/job"
	"github.com/etum-dev/WebZR/pkg/scan"
	"github.com/etum-dev/WebZR/pkg/utils"
)

// enqueueInputs feeds domains from flags, extra args, and stdin into the handler.
func enqueueInputs(flags *config.Flags, extraArgs []string, h *scanjob.Handler) error {
	var firstErr error

	enqueue := func(domain string) {
		h.AddJob(scanjob.Job{Domain: domain, Flags: flags})
	}

	if flags.DomainInput != "" {
		if config.IsFile(flags.DomainInput) {
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

	if flags.OutputFile != "" {
		if _, err := os.Stat(flags.OutputFile); err == nil {
			var check string
			fmt.Println(flags.OutputFile, "file exists, continue?")
			fmt.Scanln(&check)
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
	flags := config.FlagParse()
	scan.SetVerbose(flags.Verbose)

	workerCount := runtime.NumCPU()
	if workerCount < 2 {
		workerCount = 2
	}

	h := scanjob.NewHandler(workerCount, scanjob.RunScan)

	listener := make(chan scanjob.JobResult, workerCount)
	go h.Run(listener)

	resultBatches := make(chan []utils.ScanResult, workerCount)
	writerDone := utils.StreamResults(flags.OutputFile, resultBatches)

	// CTRL+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Collect results from workers and forward to the file writer
	collectorDone := make(chan struct{})
	go func() {
		for res := range listener {
			if res.Err != nil {
				fmt.Fprintf(os.Stderr, "scan error for %s: %v\n", res.Job.Domain, res.Err)
				continue
			}
			if len(res.Results) > 0 {
				resultBatches <- res.Results
			}
		}
		close(resultBatches)
		close(collectorDone)
	}()

	// Feed all inputs, then wait for workers or interrupt
	inputDone := make(chan error, 1)
	go func() {
		inputDone <- enqueueInputs(flags, flag.Args(), &h)
	}()

	select {
	case <-sigChan:
		fmt.Fprintf(os.Stderr, "\n\n(｡•́︿•̀｡) CTRL+C detected! Saving results and shutting down gracefully...\n")
	case err := <-inputDone:
		if err != nil {
			fmt.Fprintf(os.Stderr, "input warning: %v\n", err)
		}
	}

	// Wait for all to finish, then shut down the listener
	h.Wait()
	close(listener)
	<-collectorDone

	total := <-writerDone
	// add other message if 0 results
	fmt.Printf("\n(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Scan complete! %d result(s) written to %s\n", total, flags.OutputFile)
}
