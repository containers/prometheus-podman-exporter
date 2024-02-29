package pdcs_test

import (
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pdcs/Network", func() {
	It("Networks", func() {
		testNetwork := "exp_pdcs_test_network01"

		_, err := exec.Command("podman", "network", "create", testNetwork).Output()
		Expect(err).To(BeNil())

		podmanNetworks, err := pdcs.Networks()
		Expect(err).To(BeNil())

		netFound := false
		for _, network := range podmanNetworks {
			if network.Name == testNetwork {
				netFound = true

				break
			}
		}

		Expect(netFound).To(BeTrue())
	})
})
