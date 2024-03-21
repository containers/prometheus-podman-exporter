package e2e_em_test

import (
	"fmt"

	"github.com/containers/prometheus-podman-exporter/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pod", func() {
	It("pod metrics", func() {
		testPod01Name := "exp_e2e_test_pod01"

		pod01Inspect, err := utils.PodInformation(testPod01Name)
		Expect(err).To(BeNil())

		response := queryEndPoint()

		expectedPod01Info := fmt.Sprintf("podman_pod_info{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID, pod01Inspect.InfraID, pod01Inspect.Name)

		expectedPod01State := fmt.Sprintf("podman_pod_state{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID, pod01Inspect.InfraID, pod01Inspect.Name)

		expectedPod01Created := fmt.Sprintf("podman_pod_created_seconds{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID, pod01Inspect.InfraID, pod01Inspect.Name)

		expectedPod01Containers := fmt.Sprintf("podman_pod_containers{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID, pod01Inspect.InfraID, pod01Inspect.Name)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Info)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01State)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Created)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Containers)))
	})
})
