package e2e_test

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
		testPod02Name := "exp_e2e_test_pod02"
		testPod02ContainerName := "exp_e2e_test_pod02_cnt01"
		testPod02ContainerImage := "quay.io/quay/busybox"

		_, err := exec.Command("podman", "pod", "create", testPod01Name).Output()
		Expect(err).To(BeNil())

		_, err = exec.Command("podman", "pod", "create", testPod02Name).Output()
		Expect(err).To(BeNil())

		_, err = exec.Command("podman", "container", "create", "--pod", testPod02Name, "--name", testPod02ContainerName, testPod02ContainerImage).Output()
		Expect(err).To(BeNil())

		var (
			pod01Inspect entities.PodInspectReport
			pod02Inspect entities.PodInspectReport
		)

		pod01InspectOutput, err := exec.Command("podman", "pod", "inspect", testPod01Name).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(pod01InspectOutput, &pod01Inspect)
		Expect(err).To(BeNil())

		pod02InspectOutput, err := exec.Command("podman", "pod", "inspect", testPod02Name).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(pod02InspectOutput, &pod02Inspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()

		// podman_pod_state
		expectedPod01State := fmt.Sprintf("podman_pod_state{id=\"%s\"} 0", pod01Inspect.ID[0:12])
		expectedPod02State := fmt.Sprintf("podman_pod_state{id=\"%s\"} 0", pod02Inspect.ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01State)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02State)))

		// podman_pod_created_seconds
		expectedPod01Created := fmt.Sprintf("podman_pod_created_seconds{id=\"%s\"}", pod01Inspect.ID[0:12])
		expectedPod02Created := fmt.Sprintf("podman_pod_created_seconds{id=\"%s\"}", pod02Inspect.ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Created)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02Created)))

		// podman_pod_info
		expectedPod01Info := fmt.Sprintf("podman_pod_info{id=\"%s\",infra_id=\"%s\",name=\"%s\"} 1",
			pod01Inspect.ID[0:12], pod01Inspect.InfraContainerID[0:12], testPod01Name)
		expectedPod02Info := fmt.Sprintf("podman_pod_info{id=\"%s\",infra_id=\"%s\",name=\"%s\"} 1",
			pod02Inspect.ID[0:12], pod02Inspect.InfraContainerID[0:12], testPod02Name)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Info)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02Info)))

		// podman_pod_containers
		expectedPod01Containers := fmt.Sprintf("podman_pod_containers{id=\"%s\"} 1", pod01Inspect.ID[0:12])
		expectedPod02Containers := fmt.Sprintf("podman_pod_containers{id=\"%s\"} 2", pod02Inspect.ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Containers)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02Containers)))
	})
})
