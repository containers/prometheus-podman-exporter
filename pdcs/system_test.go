package pdcs_test

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pdcs/System", func() {
	It("SystemInfo", func() {
		system, err := pdcs.SystemInfo()
		Expect(err).To(BeNil())

		Expect(system.Buildah).ToNot(BeEmpty())
		Expect(system.Conmon).ToNot(BeEmpty())
		Expect(system.Podman).ToNot(BeEmpty())
		Expect(system.Runtime).ToNot(BeEmpty())
	})
})
