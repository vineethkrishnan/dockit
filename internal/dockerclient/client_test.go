package dockerclient

import (
	"context"
	"testing"
)

func TestNewClient_NoDaemon(t *testing.T) {
	// this allows the go test to run in CI even when a docker daemon is not available.
	// We only run this to verify the code can be executed.
	t.Log("Verifying client creation...")

	// Create a client with no options to test the function signature.
	client, err := NewClient()

	if err != nil {
		t.Logf("Expected client initialization to succeed (but ping may fail): %v", err)
	}

	if client != nil {
		err = client.Ping(context.Background())
		if err != nil {
			t.Logf("Ping correctly failed when daemon is missing: %v", err)
		}
	}
}
