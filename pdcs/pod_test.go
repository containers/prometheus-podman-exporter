package pdcs_test

import (
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pdcs/Pod", func() {
	It("Pods", func() {
		testPodName := "exp_pdcs_test_pod01"
		_, err := exec.Command("podman", "pod", "create", testPodName).Output()
		Expect(err).To(BeNil())

		podmanPods, err := pdcs.Pods()
		Expect(err).To(BeNil())

		podFound := false
		for _, pod := range podmanPods {
			if pod.Name == testPodName {
				podFound = true

				break
			}
		}

		Expect(podFound).To(BeTrue())
	})
})
