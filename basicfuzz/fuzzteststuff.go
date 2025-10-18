package basicfuzz

import (
	"log"
)


func SimpleFuzz() {
	// Create a new fuzzer targeting a websocket endpoint
	target := "ws://localhost:8000/ws"
	fuzzer := NewWSFuzzer(target)

	// Configure fuzzer settings
	fuzzer.MaxPayloads = 50 // Limit the number of payloads

	// Run the fuzzing campaign
	if err := fuzzer.Fuzz(); err != nil {
		log.Fatalf("Fuzzing failed: %v", err)
	}

	// Analyze specific results
	for _, result := range fuzzer.Results {
		if result.Anomaly {
			log.Printf("Anomaly found: Payload=%q, Response=%q\n",
				result.Payload, result.Response)
		}
	}
}

func customHeaderExample() {
	// Fuzzing with custom authentication headers
	target := "wss://example.com/secure-ws"
	fuzzer := NewWSFuzzer(target)

	// Set custom headers (e.g., for authentication)
	fuzzer.SetCustomHeaders(map[string]string{
		"Authorization": "Bearer your-token-here",
		"X-API-Key":     "your-api-key",
		"Origin":        "https://example.com",
	})

	if err := fuzzer.Fuzz(); err != nil {
		log.Fatalf("Fuzzing failed: %v", err)
	}
}

func mutationFuzzExample() {
	// Mutation-based fuzzing with seed payloads
	target := "ws://localhost:8080/ws"
	fuzzer := NewWSFuzzer(target)

	// Define seed payloads based on what you know about the API
	seedPayloads := []string{
		`{"action": "subscribe", "channel": "updates"}`,
		`{"action": "unsubscribe", "channel": "updates"}`,
		`{"action": "message", "content": "hello"}`,
		`{"action": "ping"}`,
	}

	// Run mutation-based fuzzing
	if err := fuzzer.FuzzWithMutations(seedPayloads); err != nil {
		log.Fatalf("Mutation fuzzing failed: %v", err)
	}

	// Check for crashes or interesting responses
	for _, result := range fuzzer.Results {
		if result.StatusCode >= 500 {
			log.Printf("Server error detected: %d\n", result.StatusCode)
		}
	}
}
