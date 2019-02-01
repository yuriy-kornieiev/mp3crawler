package main

import (
	"testing"
)

func TestGetEnvironment(t *testing.T) {

	cpl := CommandLineParser{}
	args := []string{"yuriyk", "crawler"}

	res, _ := cpl.GetEnvironment(args)
	if res != "yuriyk" {
		t.Error("Expected yuriyk, got ", res)
	}
}

func TestGetSource(t *testing.T) {
	cpl := CommandLineParser{}
	args := []string{"yuriyk", "crawler"}

	res, _ := cpl.GetSource(args)
	if res != "crawler" {
		t.Error("Expected crawler, got ", res)
	}
}
