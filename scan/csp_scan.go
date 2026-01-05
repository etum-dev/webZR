package scan

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/etum-dev/WebZR/utils"
)

var (
	cspHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	// Regex to extract ws:// or wss:// URLs from CSP header
	cspRegex = regexp.MustCompile(`(wss?://[^\s;"',]+)`)
)

// ScanCSP checks Content-Security-Policy headers for WebSocket URLs
// Extracts ws:// and wss:// URLs and tests them
func ScanCSP(domain string) []utils.ScanResult {
	fmt.Printf("\n(⌐■_■) Checking CSP headers for: %s\n", domain)

	// Ensure domain has http/https protocol
	url := domain
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + domain
	}

	resp, err := cspHTTPClient.Get(url)
	if err != nil {
		if verbose {
			fmt.Printf("[v] CSP request failed: %v\n", err)
		}
		return nil
	}
	defer resp.Body.Close()

	// Check for CSP header
	csp := resp.Header.Get("Content-Security-Policy")
	if csp == "" {
		if verbose {
			fmt.Printf("[v] No CSP header found\n")
		}
		return nil
	}

	// Extract all ws:// and wss:// URLs from CSP
	matches := cspRegex.FindAllString(csp, -1)
	if len(matches) == 0 {
		if verbose {
			fmt.Printf("[v] No WebSocket URLs found in CSP\n")
		}
		return nil
	}

	fmt.Printf("(ﾉ◕ヮ◕)ﾉ*:･ﾟ✧ Found %d WebSocket URL(s) in CSP header\n", len(matches))
	if len(matches) >= 1 {
		fmt.Printf("%s", matches)
	}

	var results []utils.ScanResult

	// Add all CSP findings to results
	for _, wsURL := range matches {
		fmt.Printf("  → %s\n", wsURL)

		// Determine scheme from URL
		scheme := "wss"
		if strings.HasPrefix(wsURL, "ws://") {
			scheme = "ws"
		}

		// Extract host from URL
		cleanURL := strings.TrimPrefix(wsURL, "wss://")
		cleanURL = strings.TrimPrefix(cleanURL, "ws://")

		// Create result for CSP finding
		results = append(results, utils.ScanResult{
			StatusCode: 0, // Not tested via connection
			URL:        wsURL,
			Host:       cleanURL,
			Scheme:     scheme,
			Success:    true, // Found in CSP = success
			Insecure:   scheme == "ws",
		})
	}

	return results
}
