package main

import (
	"testing"
	//"github.com/etum-dev/WebZR"
	"github.com/etum-dev/WebZR/basicfuzz"
)

func TestSimpleFuzz(t *testing.T) {
	basicfuzz.SimpleFuzz()
}

func TestHomeServer(t *testing.T) {
	basicfuzz.ServeWss()
}