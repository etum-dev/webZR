package scan

import (
	"fmt"

	"github.com/etum-dev/WebZR/utils"
)

// ScanEndpoint tries to connect to common WebSocket endpoints for a given domain
// Uses the endpoint wordlist and tests each target sequentially
func ScanEndpoint(domain string) []utils.ScanResult {
	fmt.Printf("\n(｡´-ω･)ン? Scanning endpoints for: %s\n", domain)

	targets := GenerateEndpointTargets(domain)
	if len(targets) == 0 {
		return nil
	}

	var results []utils.ScanResult

	// Test each target sequentially
	// (parallelism is handled at domain level by worker pool)
	for _, target := range targets {
		result := TestTarget(target)
		if result != nil {
			results = append(results, *result)

			// Stop on first success (optional - remove if you want all endpoints)
			break
		}
	}

	return results
}
