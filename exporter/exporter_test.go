package exporter

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestParseContainerStatsTimeout(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    time.Duration
		wantErr error
	}{
		{name: "default", value: "1s", want: time.Second},
		{name: "custom", value: "5s", want: 5 * time.Second},
		{name: "zero", value: "0s", wantErr: errInvalidContainerStatsTimeout},
		{name: "negative", value: "-1s", wantErr: errInvalidContainerStatsTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().Duration("collector.container-stats-timeout", time.Second, "")

			if err := cmd.Flags().Set("collector.container-stats-timeout", tt.value); err != nil {
				t.Fatal(err)
			}

			got, err := parseContainerStatsTimeout(cmd)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if got != tt.want {
				t.Fatalf("expected timeout %s, got %s", tt.want, got)
			}
		})
	}
}
