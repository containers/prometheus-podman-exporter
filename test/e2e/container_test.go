package e2e_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/podman/v4/pkg/domain/entities"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container", func() {
	It("container metrics", func() {
		testCnt01PodName := "exp_e2e_test_cnt01_pod01"
		testCnt01Name := "exp_e2e_test_cnt01"
		testCnt02Name := "exp_e2e_test_cnt02"
		testBusyBoxImage := "quay.io/quay/busybox:latest"

		_, err := exec.Command("podman", "pod", "create", testCnt01PodName).Output()
		Expect(err).To(BeNil())

		_, err = exec.Command("podman", "container", "create", "--pod", testCnt01PodName, "--name", testCnt01Name, testBusyBoxImage).Output()
		Expect(err).To(BeNil())

		_, err = exec.Command("podman", "container", "create", "--name", testCnt02Name, testBusyBoxImage).Output()
		Expect(err).To(BeNil())

		var (
			cnt01Inpect       []entities.ContainerInspectReport
			cnt02Inpect       []entities.ContainerInspectReport
			cnt01Pod01Inspect entities.PodInspectReport
		)

		cnt01InspectOutput, err := exec.Command("podman", "container", "inspect", testCnt01Name).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(cnt01InspectOutput, &cnt01Inpect)
		Expect(err).To(BeNil())

		cnt02InspectOutput, err := exec.Command("podman", "container", "inspect", testCnt02Name).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(cnt02InspectOutput, &cnt02Inpect)
		Expect(err).To(BeNil())

		pod01InspectOutput, err := exec.Command("podman", "pod", "inspect", testCnt01PodName).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(pod01InspectOutput, &cnt01Pod01Inspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()

		// podman_container_state
		expectedCnt01State := fmt.Sprintf("podman_container_state{id=\"%s\",pod_id=\"%s\",pod_name=\"%s\"} 0",
			cnt01Inpect[0].ID[0:12], cnt01Pod01Inspect.ID[0:12], cnt01Pod01Inspect.Name)
		expectedCnt02State := fmt.Sprintf("podman_container_state{id=\"%s\",pod_id=\"\",pod_name=\"\"} 0", cnt02Inpect[0].ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01State)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt02State)))

		// podman_container_created_seconds
		expectedCnt01Created := fmt.Sprintf("podman_container_created_seconds{id=\"%s\",pod_id=\"%s\",pod_name=\"%s\"}",
			cnt01Inpect[0].ID[0:12], cnt01Pod01Inspect.ID[0:12], cnt01Pod01Inspect.Name)
		expectedCnt02Created := fmt.Sprintf("podman_container_created_seconds{id=\"%s\",pod_id=\"\",pod_name=\"\"}", cnt02Inpect[0].ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01Created)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt02Created)))

		// podman_container_exited_seconds
		expectedCnt01ExitedSeconds := fmt.Sprintf("podman_container_exited_seconds{id=\"%s\",pod_id=\"%s\",pod_name=\"%s\"}",
			cnt01Inpect[0].ID[0:12], cnt01Pod01Inspect.ID[0:12], cnt01Pod01Inspect.Name)
		expectedCnt02ExitedSeconds := fmt.Sprintf("podman_container_exited_seconds{id=\"%s\",pod_id=\"\",pod_name=\"\"}", cnt02Inpect[0].ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01ExitedSeconds)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt02ExitedSeconds)))

		// podman_container_exit_code
		expectedCnt01ExitedCode := fmt.Sprintf("podman_container_exit_code{id=\"%s\",pod_id=\"%s\",pod_name=\"%s\"}",
			cnt01Inpect[0].ID[0:12], cnt01Pod01Inspect.ID[0:12], cnt01Pod01Inspect.Name)
		expectedCnt02ExitedCode := fmt.Sprintf("podman_container_exit_code{id=\"%s\",pod_id=\"\",pod_name=\"\"}", cnt02Inpect[0].ID[0:12])

		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01ExitedCode)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt02ExitedCode)))

		// podman_container_info
		expectedCnt01Info := fmt.Sprintf("podman_container_info{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID[0:12], cnt01Pod01Inspect.Name)
		expectedCnt02Info := fmt.Sprintf("podman_container_info{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"\",pod_name=\"\",ports=\"\"}",
			cnt02Inpect[0].ID[0:12], testBusyBoxImage, testCnt02Name)

		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01Info)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt02Info)))

	})
})
