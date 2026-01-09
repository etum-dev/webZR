package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// https://github.com/ffuf/ffuf/blob/57da720af7d1b66066cbbde685b49948f886b29c/pkg/output/stdout.go
type ScanResult struct {
	StatusCode int    `json:"status"`
	URL        string `json:"url"` //TODO: url is the found ws url(s), while host is the target scanned.
	Host       string `json:"host"`
	Scheme     string `json:"scheme"` // ws or wss
	Success    bool   `json:"success"`
	Insecure   bool   `json:"insecure"` // true if ws://
}

type VulncheckResult struct {
}

type ScanOutput struct {
	Results []ScanResult `json:"results"`
}

// TODO: Also log stuff found from CSP / JS.

// writes scan results to a JSON file
func WriteShit(filename string, results *ScanResult /*custom filename*/) error {

	output := ScanOutput{
		Results: []ScanResult{*results},
	}
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// if file exists, check if ok
	if _, err := os.Stat(filename); err == nil {
		var check string
		fmt.Println(filename, "file exists, continue?")
		fmt.Scanln(&check)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("[+] Results written to %s\n",
		filename)

	return nil
}

// writes multiple scan results to a JSON file
func WriteMultipleResults(filename string, results []ScanResult) error {
	output := ScanOutput{
		Results: results,
	}
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("[+] %d results written to %s\n", len(results), filename)

	return nil
}

// StreamResults incrementally writes results from a channel to disk and returns
// a channel that reports the total number of processed results.
func StreamResults(filename string, input <-chan []ScanResult) <-chan int {
	done := make(chan int, 1)

	go func() {
		var (
			resultFile     *os.File
			streamDisabled bool
			firstEntry     = true
			written        int
			total          int
			prefix         string
			suffix         string
		)

		defer func() {
			if resultFile != nil && !streamDisabled {
				fmt.Fprintf(resultFile, "\n%s\n", suffix)
				resultFile.Close()
				if written > 0 {
					fmt.Printf("[+] %d results written to %s\n", written, filename)
				} else {
					os.Remove(filename)
				}
			}
			done <- total
		}()

		for scanResults := range input {
			for _, res := range scanResults {
				total++
				if streamDisabled {
					continue
				}

				if resultFile == nil {
					file, err := os.Create(filename)
					if err != nil {
						fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", filename, err)
						streamDisabled = true
						continue
					}
					resultFile = file
					prefix, suffix = scanOutputDelimiters()
					fmt.Fprint(resultFile, prefix)
				}

				data, err := json.Marshal(res)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to marshal result: %v\n", err)
					continue
				}

				if firstEntry {
					fmt.Fprint(resultFile, "\n")
				} else {
					fmt.Fprint(resultFile, ",\n")
				}

				if _, err := resultFile.Write(data); err != nil {
					fmt.Fprintf(os.Stderr, "failed to write result: %v\n", err)
					resultFile.Close()
					os.Remove(filename)
					streamDisabled = true
					resultFile = nil
					continue
				}

				firstEntry = false
				written++
			}
		}
	}()

	return done
}

func scanOutputDelimiters() (string, string) {
	empty := ScanOutput{Results: []ScanResult{}}
	data, err := json.Marshal(empty)
	if err != nil {
		return "{\"results\":[", "]}"
	}

	output := string(data)
	placeholder := "[]"
	idx := strings.Index(output, placeholder)
	if idx == -1 {
		return "{\"results\":[", "]}"
	}

	prefix := output[:idx+1]
	suffix := output[idx+1:]
	return prefix, suffix
}
