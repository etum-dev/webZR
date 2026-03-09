package flaginput

import (
	"flag"
	"fmt"
	"os"

	"github.com/etum-dev/WebZR/pkg/utils"
)

type Flags struct {
	DomainInput      string // -d : can be a single domain or path to a domain list file
	OutputFile       string
	WordlistFile     string // -w : wordlist file
	Debug            bool   // -debug : enable debug
	Verbose          bool   // -v : verbose output
	Mode             string
	SubdomainMode    string // -subdomain : subdomain scanning mode (off, basic, aggressive)
	SubdomainMax     int    // -subdomain-max : maximum subdomains to test
	SubdomainWorkers int    // -subdomain-workers : number of concurrent workers for subdomain testing
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

// TODO: Add a typo/validity check. Best regards, i typo aggressive too much
func FlagParse() *Flags {
	f := &Flags{}
	flag.StringVar(&f.DomainInput, "d", "", "Single domain or path to domain list file")
	flag.StringVar(&f.OutputFile, "of", utils.OutputFileName(), "path or filename for outfile (default: scan_result(timestamp).json)")
	flag.StringVar(&f.WordlistFile, "w", "", "Path to wordlist file")
	flag.StringVar(&f.Mode, "m", "", "Mode")
	flag.BoolVar(&f.Debug, "debug", false, "Enable debug output")
	flag.BoolVar(&f.Verbose, "v", false, "Enable verbose output")
	flag.StringVar(&f.SubdomainMode, "subdomain", "off", "Subdomain scanning mode: off, basic, aggressive")
	flag.IntVar(&f.SubdomainMax, "subdomain-max", 50, "Maximum number of subdomains to test")
	flag.IntVar(&f.SubdomainWorkers, "subdomain-workers", 8, "Number of concurrent subdomain workers")

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
