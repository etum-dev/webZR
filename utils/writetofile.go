package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// https://github.com/ffuf/ffuf/blob/57da720af7d1b66066cbbde685b49948f886b29c/pkg/output/stdout.go
type ScanResult struct {
	StatusCode int    `json:"status"`
	URL        string `json:"url"`
	Host       string `json:"host"`
	Scheme     string `json:"scheme"` // ws or wss
	Success    bool   `json:"success"`
	Insecure   bool   `json:"insecure"` // true if ws://
}

type ScanOutput struct {
	Results []ScanResult `json:"results"`
}

// writes scan results to a JSON file
func WriteShit(filename string, results *ScanResult) error {

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
