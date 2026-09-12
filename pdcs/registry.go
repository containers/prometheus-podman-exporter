package pdcs

import (
	"io"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.podman.io/podman/v6/cmd/podman/registry"
)

// SetupRegistry will initialize podman registry.
func SetupRegistry() {
	// disable logrus output
	logrus.SetOutput(io.Discard)

	registry.PodmanConfig()

	_, err := registry.NewContainerEngine(&cobra.Command{}, []string{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = registry.NewImageEngine(&cobra.Command{}, []string{})
	if err != nil {
		log.Fatal(err)
	}

	cntSizeCache.cache = make(map[string]containerSize)
}
