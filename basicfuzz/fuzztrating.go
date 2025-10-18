package basicfuzz

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// FuzzResult represents the outcome of a single fuzz attempt
type FuzzResult struct {
	Payload      string
	StatusCode   int
	Response     string
	Error        error
	ResponseTime time.Duration
	Crashed      bool
	Anomaly      bool
}

// WSFuzzer handles websocket fuzzing operations
type WSFuzzer struct {
	Target       string
	Dialer       *websocket.Dialer
	MaxPayloads  int
	Timeout      time.Duration
	Results      []FuzzResult
	CustomHeader http.Header
}

// fuzzer with default settings
func NewWSFuzzer(target string) *WSFuzzer {
	return &WSFuzzer{
		Target:      target,
		MaxPayloads: 100,
		Timeout:     10 * time.Second,
		Results:     make([]FuzzResult, 0),
		Dialer: &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 10 * time.Second,
			TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
		},
		CustomHeader: http.Header{},
	}
}
func (f *WSFuzzer) SetCustomHeaders(headers map[string]string) {
	for key, value := range headers {
		f.CustomHeader.Set(key, value)
	}
}

func (f *WSFuzzer) ValidateTarget() error {
	u, err := url.Parse(f.Target)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return fmt.Errorf("invalid scheme: must be ws:// or wss://")
	}
	return nil
}

// GeneratePayloads creates various payloads for fuzzing
func (f *WSFuzzer) GeneratePayloads() []string {
	payloads := []string{
		// Basic payloads
		"test",
		"",
		" ",

		// JSON payloads
		`{"test": "value"}`,
		`{"action": "ping"}`,
		`{"command": "list"}`,

		// Malformed JSON
		`{`,
		`{"unclosed": "bracket"`,
		`{"invalid": }`,

		// Large payloads
		string(make([]byte, 1024)),
		string(make([]byte, 65536)),

		// Special characters
		"\x00\x01\x02\x03",
		"<script>alert(1)</script>",
		"'; DROP TABLE users--",
		"../../../etc/passwd",
		"%00",
		"\r\n\r\n",

		// Unicode and encoding
		"K�",
		"B5AB",
		"=%=�=�",

		// Repeated patterns
		string([]byte{0x41, 0x41, 0x41, 0x41}), // AAAA
		string([]byte{0xFF, 0xFF, 0xFF, 0xFF}),
	}

	// Add randomly generated payloads
	for i := 0; i < 20; i++ {
		payloads = append(payloads, f.generateRandomPayload())
	}

	return payloads
}

func (f *WSFuzzer) generateRandomPayload() string {
	length := rand.Intn(256) + 1
	payload := make([]byte, length)
	rand.Read(payload)
	return string(payload)
}

// MutatePayload applies mutations to a given payload
func (f *WSFuzzer) MutatePayload(original string) []string {
	mutations := []string{}

	if len(original) == 0 {
		return mutations
	}

	// Bit flip
	for i := 0; i < len(original); i++ {
		mutated := []byte(original)
		mutated[i] ^= 0xFF
		mutations = append(mutations, string(mutated))
	}

	// Byte insertion
	if len(original) < 1000 {
		for i := 0; i <= len(original); i++ {
			mutated := original[:i] + "\x00" + original[i:]
			mutations = append(mutations, mutated)
		}
	}

	// Byte deletion
	for i := 0; i < len(original); i++ {
		mutated := original[:i] + original[i+1:]
		mutations = append(mutations, mutated)
	}

	// Duplication
	mutations = append(mutations, original+original)

	// Truncation
	if len(original) > 1 {
		mutations = append(mutations, original[:len(original)/2])
	}

	return mutations
}

// SendPayload sends a single payload to the websocket and records the result
func (f *WSFuzzer) SendPayload(payload string) FuzzResult {
	result := FuzzResult{
		Payload: payload,
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), f.Timeout)
	defer cancel()

	// Connect to websocket
	conn, resp, err := f.Dialer.DialContext(ctx, f.Target, f.CustomHeader)
	if err != nil {
		result.Error = err
		result.ResponseTime = time.Since(start)
		if resp != nil {
			result.StatusCode = resp.StatusCode
		}
		return result
	}
	defer conn.Close()

	if resp != nil {
		result.StatusCode = resp.StatusCode
	}

	// Send the payload
	err = conn.WriteMessage(websocket.TextMessage, []byte(payload))
	if err != nil {
		result.Error = err
		result.ResponseTime = time.Since(start)
		return result
	}

	// Try to read response with timeout
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	msgType, message, err := conn.ReadMessage()
	if err != nil {
		// Timeout is expected for some payloads
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			result.Error = err
		}
	} else {
		result.Response = string(message)

		// Detect anomalies
		if msgType == websocket.BinaryMessage && len(message) > 10000 {
			result.Anomaly = true
		}
		if len(message) == 0 {
			result.Anomaly = true
		}
	}

	result.ResponseTime = time.Since(start)

	// Close connection gracefully
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	return result
}

// Fuzz executes the fuzzing campaign
func (f *WSFuzzer) Fuzz() error {
	if err := f.ValidateTarget(); err != nil {
		return err
	}

	log.Printf("[*] Starting fuzzing campaign against %s\n", f.Target)
	log.Printf("[*] Timeout: %v, Max Payloads: %d\n", f.Timeout, f.MaxPayloads)

	payloads := f.GeneratePayloads()

	// Limit to MaxPayloads
	if len(payloads) > f.MaxPayloads {
		payloads = payloads[:f.MaxPayloads]
	}

	log.Printf("[*] Generated %d payloads\n", len(payloads))

	for i, payload := range payloads {
		result := f.SendPayload(payload)
		f.Results = append(f.Results, result)

		if i%10 == 0 {
			log.Printf("[*] Progress: %d/%d payloads tested\n", i+1, len(payloads))
		}

		// Log interesting results
		if result.Anomaly {
			log.Printf("[!] ANOMALY detected with payload: %q\n", truncateString(payload, 50))
		}
		if result.Error != nil && result.StatusCode >= 500 {
			log.Printf("[!] Server error (%d) with payload: %q\n", result.StatusCode, truncateString(payload, 50))
		}

		// Small delay to avoid overwhelming the server
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("[+] Fuzzing complete! Tested %d payloads\n", len(f.Results))
	f.PrintSummary()

	return nil
}

// PrintSummary prints a summary of the fuzzing results
func (f *WSFuzzer) PrintSummary() {
	total := len(f.Results)
	errors := 0
	anomalies := 0
	successful := 0

	for _, result := range f.Results {
		if result.Error != nil {
			errors++
		}
		if result.Anomaly {
			anomalies++
		}
		if result.Error == nil && result.StatusCode == 101 {
			successful++
		}
	}

	log.Printf("\n=== Fuzzing Summary ===\n")
	log.Printf("Total payloads: %d\n", total)
	log.Printf("Successful: %d\n", successful)
	log.Printf("Errors: %d\n", errors)
	log.Printf("Anomalies detected: %d\n", anomalies)
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// FuzzWithMutations performs fuzzing with payload mutations
func (f *WSFuzzer) FuzzWithMutations(seedPayloads []string) error {
	if err := f.ValidateTarget(); err != nil {
		return err
	}

	log.Printf("[*] Starting mutation-based fuzzing against %s\n", f.Target)

	allPayloads := make([]string, 0)

	// Generate mutations for each seed
	for _, seed := range seedPayloads {
		mutations := f.MutatePayload(seed)
		allPayloads = append(allPayloads, mutations...)

		if len(allPayloads) >= f.MaxPayloads {
			break
		}
	}

	// Limit to MaxPayloads
	if len(allPayloads) > f.MaxPayloads {
		allPayloads = allPayloads[:f.MaxPayloads]
	}

	log.Printf("[*] Generated %d mutated payloads from %d seeds\n", len(allPayloads), len(seedPayloads))

	for i, payload := range allPayloads {
		result := f.SendPayload(payload)
		f.Results = append(f.Results, result)

		if i%10 == 0 {
			log.Printf("[*] Progress: %d/%d payloads tested\n", i+1, len(allPayloads))
		}

		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("[+] Mutation fuzzing complete!\n")
	f.PrintSummary()

	return nil
}
