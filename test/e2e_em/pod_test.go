package e2e_em_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/podman/v4/pkg/domain/entities"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pod", func() {
	It("pod metrics", func() {
		testPod01Name := "exp_e2e_test_pod01"

		var pod01Inspect entities.PodInspectReport

		pod01InspectOutput, err := exec.Command("podman", "pod", "inspect", testPod01Name).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(pod01InspectOutput, &pod01Inspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()

		expectedPod01Info := fmt.Sprintf("podman_pod_info{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID[0:12], pod01Inspect.InfraContainerID[0:12], testPod01Name)

		expectedPod01State := fmt.Sprintf("podman_pod_state{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID[0:12], pod01Inspect.InfraContainerID[0:12], testPod01Name)

		expectedPod01Created := fmt.Sprintf("podman_pod_created_seconds{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID[0:12], pod01Inspect.InfraContainerID[0:12], testPod01Name)

		expectedPod01Containers := fmt.Sprintf("podman_pod_containers{id=\"%s\",infra_id=\"%s\",name=\"%s\"}",
			pod01Inspect.ID[0:12], pod01Inspect.InfraContainerID[0:12], testPod01Name)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Info)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01State)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Created)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Containers)))
	})
})
