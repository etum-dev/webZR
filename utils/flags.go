package utils

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	DomainInput   string // -d : can be a single domain or path to a domain list file
	FuzzType      string // -fuzz : "basic", "custom", "mutation", etc.
	WordlistFile  string // -w : wordlist file
	Debug         bool   // -debug : enable debug
	WebsocketAddr string // -ws : websocket server address (if "localhost" we will spin one up)
	TestLocal     bool
}

func IsFile(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func hasStdin() bool {
	if fi, err := os.Stdin.Stat(); err == nil {
		return (fi.Mode() & os.ModeCharDevice) == 0
	}
	return false
}

func FlagParse() *Flags {
	f := &Flags{}
	flag.StringVar(&f.DomainInput, "d", "", "Single domain or path to domain list file")
	flag.StringVar(&f.FuzzType, "fuzz", "", "Fuzzing type: basic, custom, mutation")
	flag.StringVar(&f.WordlistFile, "w", "", "Path to wordlist file")
	flag.BoolVar(&f.Debug, "debug", false, "Enable debug output")
	flag.StringVar(&f.WebsocketAddr, "ws", "", "WebSocket server address")
	flag.BoolVar(&f.TestLocal, "test", false, "Use local test server")

	flag.Parse()

	if f.DomainInput == "" &&
		f.WordlistFile == "" &&
		!hasStdin() &&
		len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No input provided!")
		fmt.Fprintln(os.Stderr, "Usage: Provide -d (domain/list), -w (wordlist), piped stdin, or domain arguments")
		os.Exit(1)
	}

	return f
}
