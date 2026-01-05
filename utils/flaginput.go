package utils

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	DomainInput  string // -d : can be a single domain or path to a domain list file
	FuzzType     string // -fuzz : "lookup", "aggressive (all tests)"
	WordlistFile string // -w : wordlist file
	Debug        bool   // -debug : enable debug
	Verbose      bool   // -v : verbose output
	Mode         string
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
		return (fi.Mode() & os.ModeCharDevice) == 0 // TODO:  Unix specific - but nobody cares about windows tbh
	}
	return false
}

func FlagParse() *Flags {
	f := &Flags{}
	flag.StringVar(&f.DomainInput, "d", "", "Single domain or path to domain list file")
	flag.StringVar(&f.FuzzType, "fuzz", "", "Fuzzing type: basic, aggressive. Basic only tries to find WS on the domain. Aggressive should be used on already determined WS hosts.")
	flag.StringVar(&f.WordlistFile, "w", "", "Path to wordlist file")
	flag.StringVar(&f.Mode, "m", "", "Mode")
	flag.BoolVar(&f.Debug, "debug", false, "Enable debug output")
	flag.BoolVar(&f.Verbose, "v", false, "Enable verbose output")

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
