package pdcs

import (
	"io/ioutil"
	"log"

	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SetupRegistry will initialize podman registry.
func SetupRegistry() {
	// disable logrus output
	logrus.SetOutput(ioutil.Discard)

	registry.PodmanConfig()

	_, err := registry.NewContainerEngine(&cobra.Command{}, []string{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = registry.NewImageEngine(&cobra.Command{}, []string{})
	if err != nil {
		log.Fatal(err)
	}
}
