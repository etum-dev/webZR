package scan

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/etum-dev/WebZR/pkg/utils"
)

// SubdomainScanOptions controls subdomain scanning behavior
type SubdomainScanOptions struct {
	MaxConcurrent int           // Max concurrent subdomain tests
	Timeout       time.Duration // Timeout per subdomain test
	StopOnFirst   bool          // Stop after finding first WebSocket
	MaxSubdomains int           // Limit total subdomains to test
	PrioritizeWS  bool          // Test WebSocket-related subdomains first
}

// DefaultSubdomainOptions returns sensible defaults
func DefaultSubdomainOptions() SubdomainScanOptions {
	return SubdomainScanOptions{
		MaxConcurrent: 8,
		Timeout:       3 * time.Second,
		StopOnFirst:   false,
		MaxSubdomains: 50,
		PrioritizeWS:  true,
	}
}

func ScanSubdomain(domain string) []utils.ScanResult {
	return ScanSubdomainWithOptions(domain, DefaultSubdomainOptions())
}

// ScanSubdomainWithOptions provides configurable subdomain scanning
func ScanSubdomainWithOptions(domain string, opts SubdomainScanOptions) []utils.ScanResult {
	fmt.Printf("\nSubdomain scanning for: %s (max:%d, concurrent:%d)\n", domain, opts.MaxSubdomains, opts.MaxConcurrent)

	targets := GenerateSubdomainTargets(domain)
	if len(targets) == 0 {
		fmt.Printf(" No subdomain targets generated for %s\n", domain)
		return nil
	}

	// Limit and prioritize targets
	targets = optimizeTargetList(targets, opts)
	fmt.Printf("Testing %d optimized subdomain targets\n", len(targets))

	if len(targets) == 0 {
		return nil
	}

	// Use worker pool for concurrent subdomain testing
	return testSubdomainsConcurrently(targets, opts)
}

// optimizeTargetList prioritizes and limits subdomain targets
func optimizeTargetList(targets []Target, opts SubdomainScanOptions) []Target {
	if len(targets) == 0 {
		return targets
	}

	// Limit to max subdomains
	if len(targets) > opts.MaxSubdomains {
		targets = targets[:opts.MaxSubdomains]
	}

	return targets
}

// TODO: containsKeyword checks if URL contains a WebSocket-related keyword
func containsKeyword(url, keyword string) bool {
	if len(url) == 0 || len(keyword) == 0 {
		return false
	}

	// Convert to lowercase for case-insensitive matching
	urlLower := strings.ToLower(url)
	keywordLower := strings.ToLower(keyword)

	return strings.Contains(urlLower, keywordLower)
}

// testSubdomainsConcurrently tests subdomain targets using worker pool pattern
func testSubdomainsConcurrently(targets []Target, opts SubdomainScanOptions) []utils.ScanResult {
	if len(targets) == 0 {
		return nil
	}

	// Create channels for worker communication
	targetChan := make(chan Target, len(targets))
	resultChan := make(chan *utils.ScanResult, len(targets))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < opts.MaxConcurrent; i++ {
		wg.Add(1)
		go subdomainWorker(&wg, targetChan, resultChan, opts.Timeout)
	}

	// Send targets to workers
	go func() {
		defer close(targetChan)
		for _, target := range targets {
			targetChan <- target
		}
	}()

	// Collect results
	var results []utils.ScanResult
	var resultWg sync.WaitGroup
	resultWg.Add(1)

	go func() {
		defer resultWg.Done()
		processed := 0
		successCount := 0

		for result := range resultChan {
			processed++

			if result != nil && result.Success {
				results = append(results, *result)
				successCount++
				fmt.Printf("(｡･д･)ﾉｵ WebSocket subdomain found: %s\n", result.URL)

				// Stop on first success if configured
				if opts.StopOnFirst {
					break
				}
			}

			// Stop when all targets processed
			if processed >= len(targets) {
				break
			}
		}

		fmt.Printf("-ω- Subdomain scan completed: %d/%d successful\n", successCount, len(targets))
	}()

	// Wait for workers to finish
	wg.Wait()
	close(resultChan)

	// Wait for result collection to finish
	resultWg.Wait()

	return results
}

// subdomainWorker processes subdomain targets from the target channel
func subdomainWorker(wg *sync.WaitGroup, targets <-chan Target, results chan<- *utils.ScanResult, timeout time.Duration) {
	defer wg.Done()

	for target := range targets {
		// Test target with timeout
		result := testTargetWithTimeout(target, timeout)
		results <- result
	}
}

// testTargetWithTimeout tests a target with a specific timeout
func testTargetWithTimeout(target Target, timeout time.Duration) *utils.ScanResult {
	// Create a channel to receive the result
	resultChan := make(chan *utils.ScanResult, 1)

	// Start the test in a goroutine
	go func() {
		result := TestTarget(target)
		resultChan <- result
	}()

	// Wait for result or timeout
	select {
	case result := <-resultChan:
		return result
	case <-time.After(timeout):
		return nil // Timeout
	}
}
