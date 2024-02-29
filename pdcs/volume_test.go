package pdcs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/containers/prometheus-podman-exporter/pdcs"
)

var _ = Describe("Volume", func() {
	It("PDCS Volumes", func() {
		_, err := pdcs.Volumes()
		Expect(err).To(BeNil())
	})
})
