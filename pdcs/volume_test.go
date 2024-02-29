package pdcs_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/containers/prometheus-podman-exporter/pdcs"
)

var _ = Describe("Pdcs/Volume", func() {
	It("Volumes", func() {
		testVolName := "exp_pdcs_test_vol01"

		_, err := exec.Command("podman", "volume", "create", testVolName).Output()
		Expect(err).To(BeNil())

		podmanVolumes, err := pdcs.Volumes()
		Expect(err).To(BeNil())

		volFound := false
		for _, vol := range podmanVolumes {
			if vol.Name == testVolName {
				volFound = true

				break
			}
		}

		Expect(volFound).To(BeTrue())
	})
})
