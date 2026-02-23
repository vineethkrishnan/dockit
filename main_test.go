package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
