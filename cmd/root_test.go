package cmd

import (
	"testing"
	"time"
)

func TestContainerStatsTimeoutDefault(t *testing.T) {
	timeout, err := rootCmd.Flags().GetDuration("collector.container-stats-timeout")
	if err != nil {
		t.Fatal(err)
	}

	if timeout != time.Second {
		t.Fatalf("expected default container stats timeout to be 1s, got %s", timeout)
	}
}
