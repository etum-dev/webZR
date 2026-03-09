package scan

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/etum-dev/WebZR/pkg/utils"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

// Target represents a single URL to test for WebSocket support
type Target struct {
	URL    string
	Source string // "endpoint", "subdomain", "csp", "base"
}

// GenerateEndpointTargets creates all endpoint-based targets for a domain
// Returns a list of URLs to test (does NOT perform connections)
func GenerateEndpointTargets(domain string) []Target {
	endpoints, err := getEndpointWordlist()
	if err != nil || len(endpoints) == 0 {
		// Fallback to just the base domain
		return []Target{{URL: strings.TrimRight(domain, "/"), Source: "base"}}
	}

	base := strings.TrimRight(domain, "/")
	targets := make([]Target, 0, len(endpoints)+1)

	// Generate target for each endpoint
	for _, endpoint := range endpoints {
		targets = append(targets, Target{
			URL:    joinEndpoint(base, endpoint),
			Source: "endpoint",
		})
	}

	// Always try base domain too
	targets = append(targets, Target{
		URL:    base,
		Source: "base",
	})

	return targets
}

// Returns a list of URLs to test (does not perform connections)
// it first retrieves valid targets via subfinder,
// if none found, it will instead fall back on a default wordlist.
func GenerateSubdomainTargets(domain string) []Target {
	// first, do a subfinder on the target.
	subfinderOpts := &runner.Options{
		Threads:            10, // TODO: this should be dinamic based on host
		Timeout:            30,
		MaxEnumerationTime: 20,
	}
	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		fmt.Printf("subfinder failed: %v", err)
	}
	output := &bytes.Buffer{}
	var sourceMap map[string]map[string]struct{}
	if sourceMap, err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), domain, []io.Writer{output}); err != nil {
		fmt.Printf("Enumeration failed on domain: %v", err)
	}
	var targets []Target
	for subdomain, _ := range sourceMap {
		targets = append(targets, Target{
			URL:    subdomain,
			Source: "subdomain",
		})

		fmt.Println(targets)
	}

	// if bad results, manually read wordlist and shoot
	/*subdomains, err := getSubdomainWordlist()
	if err != nil || len(subdomains) == 0 {
		return nil
	}

	host := utils.ExtractHostname(domain)
	if host == "" {
		return nil
	}

	// Get endpoints to combine with subdomains
	endpoints, err := getEndpointWordlist()
	if err != nil || len(endpoints) == 0 {
		// Just use empty path (base domain)
		endpoints = []string{""}
	}

	var targets []Target
	for _, sub := range subdomains {
		sub = strings.TrimSpace(sub)
		if sub == "" {
			continue
		}

		fullDomain := fmt.Sprintf("%s.%s", sub, host)

		// For each subdomain, try all endpoints
		for _, endpoint := range endpoints {
			targets = append(targets, Target{
				URL:    joinEndpoint(fullDomain, endpoint),
				Source: "subdomain",
			})
		}
	}
	*/
	return targets
}

// TestTarget attempts a WebSocket connection to a single target
// This is what workers should call
func TestTarget(target Target) *utils.ScanResult {
	result := SendConnRequest(target.URL)

	// Only return successful results
	if result != nil && result.Success {
		return result
	}

	return nil
}
