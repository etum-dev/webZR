package scan

import (
	"fmt"
	"strings"

	"github.com/etum-dev/WebZR/utils"
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

// GenerateSubdomainTargets creates all subdomain-based targets for a domain
// Returns a list of URLs to test (does NOT perform connections)
func GenerateSubdomainTargets(domain string) []Target {
	subdomains, err := getSubdomainWordlist()
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
