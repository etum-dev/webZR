package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/etum-dev/WebZR/scan"
)

func TestScanEndpoint(t *testing.T) {
	urlfile, err := os.Open("../donotpushthisyouwillgetfired.txt")
	if err != nil {
		fmt.Println("xd")
	}
	defer urlfile.Close()

	scanner := bufio.NewScanner(urlfile)
	for scanner.Scan() {
		scan.ScanEndpoint(scanner.Text())
	}

}
