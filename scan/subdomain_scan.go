package scan

import (
	"fmt"

	"github.com/etum-dev/WebZR/utils"
)

// ScanSubdomain tries to find WebSocket subdomains of a domain
func ScanSubdomain(domain string) []utils.ScanResult {
	fmt.Printf("\n(っ´▽`)っ Scanning subdomains for: %s\n", domain)

	targets := GenerateSubdomainTargets(domain)
	if len(targets) == 0 {
		fmt.Printf("No subdomain targets generated for %s\n", domain)
		return nil
	}

	var results []utils.ScanResult

	// Test each target sequentially
	// The main worker pool provides parallelism across domains
	for _, target := range targets {
		result := TestTarget(target)
		if result != nil {
			results = append(results, *result)
			fmt.Printf("(((((っ･ω･)っ ﾌﾞｰﾝ Websocket subdomain found: %s\n", result.URL)

			// Optionally stop on first success
			// Remove this break to find all subdomain websockets
			break
		}
	}

	return results
}
