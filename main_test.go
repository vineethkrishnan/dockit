package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Just verify it compiles and basic imports work to appease CI for now.
	os.Exit(m.Run())
}
