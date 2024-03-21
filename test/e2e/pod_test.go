package e2e_test

import (
	"fmt"
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/test/utils"
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

		pod01Inspect, err := utils.PodInformation(testPod01Name)
		Expect(err).To(BeNil())

		pod02Inspect, err := utils.PodInformation(testPod02Name)
		Expect(err).To(BeNil())

		response := queryEndPoint()

		// podman_pod_state
		expectedPod01State := fmt.Sprintf("podman_pod_state{id=\"%s\"} 0", pod01Inspect.ID)
		expectedPod02State := fmt.Sprintf("podman_pod_state{id=\"%s\"} 0", pod02Inspect.ID)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01State)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02State)))

		// podman_pod_created_seconds
		expectedPod01Created := fmt.Sprintf("podman_pod_created_seconds{id=\"%s\"}", pod01Inspect.ID)
		expectedPod02Created := fmt.Sprintf("podman_pod_created_seconds{id=\"%s\"}", pod02Inspect.ID)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Created)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02Created)))

		// podman_pod_info
		expectedPod01Info := fmt.Sprintf("podman_pod_info{id=\"%s\",infra_id=\"%s\",name=\"%s\"} 1",
			pod01Inspect.ID, pod01Inspect.InfraID, pod01Inspect.Name)
		expectedPod02Info := fmt.Sprintf("podman_pod_info{id=\"%s\",infra_id=\"%s\",name=\"%s\"} 1",
			pod02Inspect.ID, pod02Inspect.InfraID, pod02Inspect.Name)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Info)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02Info)))

		// podman_pod_containers
		expectedPod01Containers := fmt.Sprintf("podman_pod_containers{id=\"%s\"} 1", pod01Inspect.ID)
		expectedPod02Containers := fmt.Sprintf("podman_pod_containers{id=\"%s\"} 2", pod02Inspect.ID)

		Expect(response).Should(ContainElement(ContainSubstring(expectedPod01Containers)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedPod02Containers)))
	})
})
