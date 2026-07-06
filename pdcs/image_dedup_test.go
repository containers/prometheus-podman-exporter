package pdcs_test

import (
	"testing"

	"github.com/containers/prometheus-podman-exporter/pdcs"
)

func TestDedupImages(t *testing.T) {
	// A digest-pinned image with two name references. Both decompose to the same
	// repository and "<none>" tag and carry the same image-level digest, so the
	// two Image entries are identical (the traefik case from the bug report).
	traefikRepo := "docker.io/library/traefik"
	input := []pdcs.Image{
		{ID: "2863dabaa216", Repository: traefikRepo, Tag: "<none>", Digest: "sha256:5a52522"},
		{ID: "2863dabaa216", Repository: traefikRepo, Tag: "<none>", Digest: "sha256:5a52522"},
		{ID: "e76fd176", Repository: "docker.io/grafana/grafana", Tag: "<none>", Digest: "sha256:121a7a9"},
		// same image id, different tags must remain distinct series.
		{ID: "abcabc", Repository: "docker.io/library/nginx", Tag: "1.27"},
		{ID: "abcabc", Repository: "docker.io/library/nginx", Tag: "latest"},
	}

	got := pdcs.DedupImages(input)

	if len(got) != 4 {
		t.Fatalf("expected 4 unique images, got %d: %+v", len(got), got)
	}

	traefik := 0
	nginx := 0

	for _, img := range got {
		switch img.Repository {
		case traefikRepo:
			traefik++
		case "docker.io/library/nginx":
			nginx++
		}
	}

	if traefik != 1 {
		t.Errorf("expected the duplicate traefik references to collapse to 1, got %d", traefik)
	}

	if nginx != 2 {
		t.Errorf("expected both nginx tags to survive, got %d", nginx)
	}
}
