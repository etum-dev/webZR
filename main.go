package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/etum-dev/WebZR/basicfuzz"
	"github.com/etum-dev/WebZR/scan"
)

// processIn consumes optional file, args and/or stdin and prints each line/item.
// You can replace the fmt.Println calls with whatever processing you need.
func processIn(inputFile string, args []string, hasStdin bool) {
	// process positional args first (if any)
	if len(args) > 0 {
		for _, a := range args {
			fmt.Println("arg:", a)
		}
	}

	// process file if supplied
	if inputFile != "" {
		f, err := os.Open(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open file %s: %v\n", inputFile, err)
		} else {
			defer f.Close()
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				fmt.Println("file:", sc.Text())
			}
			if err := sc.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", inputFile, err)
			}
		}
	}

	// process piped stdin if present
	if hasStdin {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			fmt.Println("stdin:", sc.Text())
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading stdin: %v\n", err)
		}
	}
}

func main() {
	// flags
	isTest := flag.Bool("test", false, "Use the local test server")
	fuzzOpt := flag.Int("fuzz", 0, "Fuzzing option. 1 for basic, 2 for custom header, 3 for mutation")
	fileOpt := flag.String("file", "", "Path to input file (optional)")
	singleDomain := flag.String("domain", "", "Single domain")

	flag.Parse()

	// detect positional args
	args := flag.Args()

	// detect whether there's piped stdin
	hasStdin := false
	if fi, err := os.Stdin.Stat(); err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			hasStdin = true
		}
	}

	// Decide whether to call processIn:
	// call it when user passed a file flag, or provided positional args, or piped stdin.
	if *fileOpt != "" || len(args) > 0 || hasStdin {
		processIn(*fileOpt, args, hasStdin)
	} else if *singleDomain != "" {
		scan.SendConnRequest(*singleDomain)
	} else {

		fmt.Println("Usage: provide -file=domains.txt, -domain=example.com, or domains as arguments")
	}

	if *isTest {
		basicfuzz.ServeWs()
	}

	switch *fuzzOpt {
	case 1:
		fmt.Println("Doing a simple fuzz.")
		basicfuzz.SimpleFuzz()
	case 0:
		fmt.Println("Not fuzzing")
	}

}
