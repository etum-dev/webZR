package scanjob

import (
	"fmt"
	"time"

	"github.com/etum-dev/WebZR/pkg/scan"
	"github.com/etum-dev/WebZR/pkg/utils"
)

// RunScan executes the full scan pipeline for a single Job.
func RunScan(job Job) JobResult {
	domain := utils.CheckDomain(job.Domain)
	if domain == "" {
		return JobResult{Job: job, Err: fmt.Errorf("empty domain")}
	}

	var results []utils.ScanResult
	mode := job.Flags.Mode
	if mode == "" {
		mode = "basic"
	}

	// Always try CSP first
	if cspResults := scan.ScanCSP(domain); len(cspResults) > 0 {
		results = append(results, cspResults...)
	}
	// crawl page JS and search for wss strings
	if jsResults := scan.JSCrawler(domain); len(jsResults) > 0 {
		results = append(results, jsResults...)
	}

	// aggressive = request heavy scans
	if mode == "aggressive" {
		if epResults := scan.ScanEndpoint(domain); len(epResults) > 0 {
			results = append(results, epResults...)
		}
	}

	// Enhanced subdomain scanning with configurable options
	subdomainMode := job.Flags.SubdomainMode
	if subdomainMode == "" {
		subdomainMode = "off"
	}

	if subdomainMode != "off" {
		var subResults []utils.ScanResult

		if subdomainMode == "basic" {
			opts := scan.SubdomainScanOptions{
				MaxConcurrent: job.Flags.SubdomainWorkers,
				Timeout:       5 * time.Second,
				StopOnFirst:   true,
				MaxSubdomains: func() int {
					if job.Flags.SubdomainMax < 25 {
						return job.Flags.SubdomainMax
					}
					return 25
				}(),
				PrioritizeWS: true,
			}
			subResults = scan.ScanSubdomainWithOptions(domain, opts)
		} else if subdomainMode == "aggressive" {
			opts := scan.SubdomainScanOptions{
				MaxConcurrent: job.Flags.SubdomainWorkers,
				Timeout:       3 * time.Second,
				StopOnFirst:   false,
				MaxSubdomains: job.Flags.SubdomainMax,
				PrioritizeWS:  true,
			}
			subResults = scan.ScanSubdomainWithOptions(domain, opts)
		}

		if len(subResults) > 0 {
			results = append(results, subResults...)
		}
	}

	return JobResult{
		Job:     job,
		Results: results,
	}
}
