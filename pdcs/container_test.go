package pdcs_test

import (
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pdcs/Container", func() {
	It("Containers", func() {
		testImage := "docker.io/library/busybox"

		_, err := exec.Command("podman", "image", "pull", testImage).Output()
		Expect(err).To(BeNil())

		testCntName := "exp_pdcs_test_container01"
		_, err = exec.Command("podman", "container", "create", "--name", testCntName, testImage).Output()
		Expect(err).To(BeNil())

		podmanContainers, err := pdcs.Containers()
		Expect(err).To(BeNil())

		cntFound := false
		for _, container := range podmanContainers {
			if container.Name == testCntName {
				cntFound = true

				break
			}
		}

		Expect(cntFound).To(BeTrue())
	})
})
