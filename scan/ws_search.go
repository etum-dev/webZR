package scan

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/websocket"

	"github.com/etum-dev/WebZR/utils"
)

func CheckJS(url string) {
	// Parse sites js to find if it has ws estab.
	// Should be delegated to tools like katana and idk rg
}

// ScanEndpoint tries to connect to common WebSocket endpoints for a given domain
// It reads endpoints from ws-endpoints.txt and attempts connections
func ScanEndpoint(url string) utils.ScanResult {
	fmt.Printf("\n[*] Scanning endpoints for: %s\n", url)

	wordlist, err := os.Open("ws-endpoints.txt")
	if err != nil {
		fmt.Printf("[-] Cannot open ws-endpoints.txt: %v\n", err)
		fmt.Println("[*] Skipping endpoint scan, trying base domain only...")

		// Try just the base domain
		result := SendConnRequest(url)
		if result != nil {
			return *result
		}

		// Return empty result if nothing worked
		return utils.ScanResult{
			Host:    url,
			Success: false,
		}
	}
	defer wordlist.Close()

	// Try each endpoint from the wordlist
	scanner := bufio.NewScanner(wordlist)
	for scanner.Scan() {
		endpoint := scanner.Text()
		fullURL := url + endpoint

		fmt.Printf("[*] Trying endpoint: %s\n", endpoint)
		result := SendConnRequest(fullURL)

		// If we found a working endpoint, return it immediately
		if result != nil && result.Success {
			return *result
		}
	}

	// No endpoints worked
	return utils.ScanResult{
		Host:    url,
		Success: false,
	}
}

func ScanSubdomain(url string) utils.ScanResult {
	retryAttempts := 3
	attempts := 0

	subdomainFile, err := os.Open("ws-subdomain.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file for subdomains: %v\n", err)
	} else {
		defer subdomainFile.Close()
		subdomainSc := bufio.NewScanner(subdomainFile)
		for subdomainSc.Scan() {
			// net.LookupNS
			// if (something fucked in request){ attempts +1 }
			if attempts == retryAttempts {
				fmt.Println("fuck off m8")
			}
		}
	}
	placeholder := &utils.ScanResult{
		StatusCode: 555,
		URL:        url,
		Host:       url,
		Scheme:     "ws",
		Success:    true,
		Insecure:   false,
	}
	return *placeholder
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
		// Build the full WebSocket URL
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

		// Close connection if established
		if conn != nil {
			defer conn.Close()
		}

		// HTTP 101 = "Switching Protocols" = WebSocket handshake successful
		if resp.StatusCode == 101 {
			isInsecure := (scheme == "ws")

			fmt.Printf("[+] WebSocket connection established: %s (Status %d)\n", wsUrl, resp.StatusCode)

			if isInsecure {
				fmt.Printf("[!] WARNING: Insecure WebSocket (ws://) connection accepted!\n")
			}

			// Return the successful result
			return &utils.ScanResult{
				StatusCode: resp.StatusCode,
				URL:        wsUrl,
				Host:       domain,
				Scheme:     scheme,
				Success:    true,
				Insecure:   isInsecure,
			}
		} else {
			fmt.Printf("[-] WebSocket upgrade failed %s: Status %d\n", wsUrl, resp.StatusCode)
		}
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

func FindMoreEndpoints() {
	// Calls functions that look for wss in resp headers, javascript

}
func CSPSearch(rhttp string) (_ []string) {
	// Check if http header has any wss domains
	r, err := http.Head(rhttp)
	if err != nil {
		fmt.Println("HTTP Request failed: ", err)
	}
	csp := r.Header.Get("content-security-policy")
	/*if err != nil {
		fmt.Println("CSP Header check failed: ", err)
	}*/

	if err != nil {
		fmt.Println("hatsune miku")
	} else {
		re := regexp.MustCompile(`ws{1,2}\:\/\/`)
		match := re.MatchString(csp)
		if match {
			wsheader := re.FindAllString(csp, -1)
			return wsheader
		}

	}
	return []string{""}

}

func CorsITaket(url string, ownserver string) {
	// Check if it validates origin'
	/* if (cors not validated) {
		send request with own server value
	}*/

}
