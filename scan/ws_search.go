package scan

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/etum-dev/WebZR/utils"
	"github.com/gorilla/websocket"
)

var verbose bool

// request/response logging.
func SetVerbose(v bool) {
	verbose = v
}

func logVerboseHandshake(wsURL string, resp *http.Response, err error) {
	if !verbose {
		return
	}

	fmt.Printf("[v] Dial attempt: %s\n", wsURL)
	if err != nil {
		fmt.Printf("[v] Dial error: %v\n", err)
	}

	if resp == nil {
		return
	}

	if resp.Request != nil {
		req := resp.Request
		fmt.Printf("[v] Request: %s %s\n", req.Method, req.URL.String())
		printVerboseHeaders(req.Header)
	}

	fmt.Printf("[v] Response: %s\n", resp.Status)
	printVerboseHeaders(resp.Header)
}

func printVerboseHeaders(headers http.Header) {
	if len(headers) == 0 {
		fmt.Println("[v]   <no headers>")
		return
	}

	for name, values := range headers {
		fmt.Printf("[v]   %s: %s\n", name, strings.Join(values, ", "))
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
		logVerboseHandshake(wsUrl, resp, err)
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

		fmt.Printf("(　´，_ゝ｀)ﾌﾟｯ WebSocket upgrade failed %s: Status %d\n", wsUrl, resp.StatusCode)
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
