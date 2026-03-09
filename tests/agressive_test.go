package tests

import (
	"testing"
	"time"

	"github.com/etum-dev/WebZR/pkg/scan"
)

func TestGenerateSubdomainTargets(t *testing.T) {
	// this test should return subfinder data
	// hardcoding a domain
	domain := "websocket.org"
	var wantsub [2]string
	wantsub[0] = "echo.websocket.org"
	wantsub[1] = "www.websocket.org"

	x := scan.GenerateSubdomainTargets(domain)
	got := make([]string, len(x))
	for i, target := range x {
		got[i] = target.URL
	}
	t.Logf("found %v, want: %v", got, wantsub)

}

func TestAgressiveModeSubdomain_Scan(t *testing.T) {
	t.Skip("integration: remove skip to run against a live target")

	domain := "example.com"
	opts := scan.SubdomainScanOptions{
		MaxConcurrent: 4,
		Timeout:       3 * time.Second,
		StopOnFirst:   false,
		MaxSubdomains: 10,
		PrioritizeWS:  true,
	}

	results := scan.ScanSubdomainWithOptions(domain, opts)
	t.Logf("found %d results for %s", len(results), domain)
}
