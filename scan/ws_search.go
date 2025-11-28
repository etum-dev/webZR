package scan

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/etum-dev/WebZR/utils"
	"github.com/gorilla/websocket"
)

var possibleWs []string

// ScanEndpoint tries to connect to common WebSocket endpoints for a given domain
// It reads endpoints from ws-endpoints.txt and attempts connections
func ScanEndpoint(url string) utils.ScanResult {
	fmt.Printf("\n(｡´-ω･)ン? Scanning endpoints for: %s\n", url)

	endpoints, err := getEndpointWordlist()
	if err != nil || len(endpoints) == 0 {
		fmt.Printf("(>^<。)グスン Cannot load ws-endpoints.txt: %v\n", err)
		fmt.Println("(。┰ω┰。) Skipping endpoint fuzz, trying base domain only...")
		if result := SendConnRequest(strings.TrimRight(url, "/")); result != nil {
			return *result
		}
		return utils.ScanResult{
			Host:    url,
			Success: false,
		}
	}

	base := strings.TrimRight(url, "/")

	for _, endpoint := range endpoints {
		fullURL := joinEndpoint(base, endpoint)
		result := SendConnRequest(fullURL)
		if result != nil && result.Success {
			return *result
		}
	}

	if result := SendConnRequest(base); result != nil {
		return *result
	}

	return utils.ScanResult{
		Host:    url,
		Success: false,
	}
}

func joinEndpoint(base, endpoint string) string {
	base = strings.TrimRight(base, "/")
	endpoint = strings.TrimSpace(endpoint)

	if endpoint == "" {
		return base
	}

	if strings.HasPrefix(endpoint, "/") || strings.HasPrefix(endpoint, "?") || strings.HasPrefix(endpoint, "&") {
		return base + endpoint
	}

	return base + "/" + endpoint
}

const subdomainConcurrency = 10

func ScanSubdomain(url string) []utils.ScanResult {
	subdomains, err := getSubdomainWordlist()
	if err != nil || len(subdomains) == 0 {
		fmt.Printf("failed to load subdomain list: %v\n", err)
		return nil
	}

	host := utils.ExtractHostname(url)
	if host == "" {
		fmt.Printf("could not derive hostname from %q, skipping subdomain scan\n", url)
		return nil
	}

	endpoints, endpointErr := getEndpointWordlist()
	if endpointErr != nil {
		fmt.Printf("failed to load endpoint list, falling back to root path only: %v\n", endpointErr)
	}
	targetEndpoints := append(make([]string, 0, len(endpoints)+1), endpoints...)
	targetEndpoints = append(targetEndpoints, "") // always probe bare host

	resultsChan := make(chan utils.ScanResult, len(subdomains))
	var wg sync.WaitGroup
	sem := make(chan struct{}, subdomainConcurrency)

	probeSubdomain := func(domain string) *utils.ScanResult {
		trimmed := strings.TrimRight(domain, "/")
		for _, endpoint := range targetEndpoints {
			fullTarget := joinEndpoint(trimmed, endpoint)
			if res := SendConnRequest(fullTarget); res != nil && res.Success {
				return res
			}
		}
		return nil
	}

	for _, subdomain := range subdomains {
		subdomain := strings.TrimSpace(subdomain)
		if subdomain == "" {
			continue
		}

		wg.Add(1)
		go func(sd string) {
			defer wg.Done()
			fullDomain := fmt.Sprintf("%s.%s", sd, host)
			//fmt.Println("trying", fullDomain)

			sem <- struct{}{}
			result := probeSubdomain(fullDomain)
			<-sem

			if result != nil && result.Success {
				resultsChan <- *result
				fmt.Println("(((((っ･ω･)っ ﾌﾞｰﾝ Websocket subdomain found:", result.URL)
			}
		}(subdomain)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []utils.ScanResult
	for res := range resultsChan {
		results = append(results, res)
	}

	return results
}

// SendConnRequest attempts to establish WebSocket connections using both wss:// and ws://
// Returns a ScanResult with connection details, or nil if no connection succeeded
func SendConnRequest(domain string) *utils.ScanResult {
	/* Attempts to connect via ws and wss. Should flag for:
	- If connection possible
	- allows for ws://
	*/

	schemes := []string{"wss", "ws"} // Try secure first, then insecure

	for _, scheme := range schemes {
		wsUrl := scheme + "://" + domain

		dialer := websocket.Dialer{
			HandshakeTimeout: 5 * time.Second,
		}

		conn, resp, err := dialer.Dial(wsUrl, nil)
		if err != nil {
			continue
		}

		if resp == nil {
			fmt.Printf("[-] %s: no response received\n", wsUrl)
			continue
		}

		if resp.Body != nil {
			resp.Body.Close()
		}

		if conn != nil {
			conn.Close()
		}

		if resp.StatusCode == http.StatusSwitchingProtocols {
			isInsecure := (scheme == "ws")

			fmt.Printf("( ｡･_･｡)人(｡･_･｡ ) WebSocket connection established: %s (Status %d)\n", wsUrl, resp.StatusCode)

			if isInsecure {
				fmt.Printf("[!](・∀・)イイ!! WARNING: Insecure WebSocket (ws://) connection accepted!\n")
			}

			return &utils.ScanResult{
				StatusCode: resp.StatusCode,
				URL:        wsUrl,
				Host:       domain,
				Scheme:     scheme,
				Success:    true,
				Insecure:   isInsecure,
			}
		}

		fmt.Printf("[-] WebSocket upgrade failed %s: Status %d\n", wsUrl, resp.StatusCode)
	}

	// No scheme worked - return failure result
	return &utils.ScanResult{
		StatusCode: 0,
		URL:        domain,
		Host:       domain,
		Scheme:     "none",
		Success:    false,
		Insecure:   false,
	}
}

//

// Reads JS discovered on page and extracts wss:// urls
func JSCrawler(url string) error {
	var html string
	// super basic, can 1000% be improved.
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	) // asså usch
	if err != nil {
		log.Fatal("Error automation logic: ", err)
	}
	fmt.Println(html)
	return nil

}

var (
	cspHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	cspRegex = regexp.MustCompile(`(ws{1,2}://[^\s;"',]+)`)
)

func CSPSearch(rhttp string) []string {
	parsedurl := utils.AppendProto(rhttp)
	if parsedurl == "" {
		return nil
	}

	resp, err := cspHTTPClient.Get(parsedurl)
	if err != nil {
		fmt.Println("Request failed:", err)
		return nil
	}
	defer resp.Body.Close()

	csp := resp.Header.Get("Content-Security-Policy")
	if csp == "" {
		return nil
	}

	matches := cspRegex.FindAllString(csp, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	candidates := make([]string, 0, len(matches))

	for _, match := range matches {
		clean := strings.TrimRight(match, ",'\";")
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		candidates = append(candidates, clean)
	}

	return candidates
}
