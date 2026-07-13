package pdcs

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.podman.io/podman/v6/libpod/define"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

func TestSetContainerStatsTimeout(t *testing.T) {
	original := getContainerStatsTimeout()
	if original != time.Second {
		t.Fatalf("expected default container stats timeout to be 1s, got %s", original)
	}

	t.Cleanup(func() {
		if err := SetContainerStatsTimeout(original); err != nil {
			t.Fatal(err)
		}
	})

	if err := SetContainerStatsTimeout(5 * time.Second); err != nil {
		t.Fatal(err)
	}

	if timeout := getContainerStatsTimeout(); timeout != 5*time.Second {
		t.Fatalf("expected container stats timeout to be 5s, got %s", timeout)
	}

	if err := SetContainerStatsTimeout(0); !errors.Is(err, errInvalidContainerStatsTimeout) {
		t.Fatalf("expected error %v, got %v", errInvalidContainerStatsTimeout, err)
	}

	if timeout := getContainerStatsTimeout(); timeout != 5*time.Second {
		t.Fatalf("invalid configuration changed container stats timeout to %s", timeout)
	}
}

func TestWaitForContainerStats(t *testing.T) {
	stats := []define.ContainerStats{{ContainerID: "test-container"}}
	statsErr := errors.New("stats failed")

	tests := []struct {
		name    string
		setup   func(context.CancelFunc, chan entities.ContainerStatsReport)
		want    []define.ContainerStats
		wantErr error
	}{
		{
			name: "report",
			setup: func(_ context.CancelFunc, reports chan entities.ContainerStatsReport) {
				reports <- entities.ContainerStatsReport{Stats: stats}
			},
			want: stats,
		},
		{
			name: "report error",
			setup: func(_ context.CancelFunc, reports chan entities.ContainerStatsReport) {
				reports <- entities.ContainerStatsReport{Error: statsErr}
			},
			wantErr: statsErr,
		},
		{
			name: "closed channel",
			setup: func(_ context.CancelFunc, reports chan entities.ContainerStatsReport) {
				close(reports)
			},
			wantErr: ErrDeadline,
		},
		{
			name: "cancelled context",
			setup: func(cancel context.CancelFunc, reports chan entities.ContainerStatsReport) {
				cancel()
				close(reports)
			},
			wantErr: ErrDeadline,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			reports := make(chan entities.ContainerStatsReport, 1)
			tt.setup(cancel, reports)

			got, err := waitForContainerStats(ctx, reports)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d stats reports, got %d", len(tt.want), len(got))
			}

			if len(got) > 0 && got[0].ContainerID != tt.want[0].ContainerID {
				t.Fatalf("expected container ID %q, got %q", tt.want[0].ContainerID, got[0].ContainerID)
			}
		})
	}
}
