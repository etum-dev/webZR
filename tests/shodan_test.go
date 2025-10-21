package main

import(
	"testing"
	"github.com/etum-dev/WebZR/scan"
)

func TestAuthShodan(t *testing.T){
	t.Parallel()
}

func TestSearchShodan(t *testing.T){
	scan.SearchShodan()
}