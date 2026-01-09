package scan

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/etum-dev/WebZR/utils"
)

// JSCrawler reads JS discovered on page and extracts wss:// urls
func JSCrawler(domain string) []utils.ScanResult {
	pdomain := utils.AppendHttpsProto(domain)

	var html string

	// super basic, can 1000% be improved.
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(pdomain),
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Print("Error automation logic: ", err, domain)

	}
	matches := doTheRegex(html)
	var results []utils.ScanResult

	// Convert found URLs into ScanResult objects
	for _, match := range matches {
		results = append(results, utils.ScanResult{
			StatusCode: 0,
			URL:        match,
			Host:       domain,
			Scheme:     "javascript",
			Success:    true,
		})
	}

	return results
}

func doTheRegex(html string) []string {
	ws_patterns := []string{
		`ws(s?):\/\/([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?`,
		`new\s+WebSocket\s*\(\s*['"]([^'"]+)['"]`,
		`WebSocket\s*\(\s*['"]([^'"]+)['"]`,
		`['"]ws(s?)://[^'"]+['"]`,
		`socket\.io[^'"]*['"]([^'"]+)['"]`,
	}

	var allMatches []string
	seen := make(map[string]bool) // Deduplicate matches

	// Compile and test each pattern
	for _, pattern := range ws_patterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("Error compiling regex pattern %s: %v", pattern, err)
			continue
		}

		matches := regex.FindAllString(html, -1)
		for _, match := range matches {
			// Clean up the match (remove quotes if present)
			cleaned := cleanMatch(match)
			if cleaned != "" && !seen[cleaned] {
				seen[cleaned] = true
				allMatches = append(allMatches, cleaned)
			}
		}
	}

	return allMatches
}

// cleanMatch removes quotes and extracts clean WebSocket URLs
func cleanMatch(match string) string {
	// Remove surrounding quotes
	cleaned := regexp.MustCompile(`^['"]|['"]$`).ReplaceAllString(match, "")

	// Extract URL from constructor patterns like: new WebSocket("ws://example.com")
	if constructorMatch := regexp.MustCompile(`(?:new\s+)?WebSocket\s*\(\s*['"]([^'"]+)['"]`).FindStringSubmatch(match); len(constructorMatch) > 1 {
		return constructorMatch[1]
	}

	return cleaned
}
