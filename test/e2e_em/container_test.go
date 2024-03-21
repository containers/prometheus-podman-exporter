package e2e_em_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/prometheus-podman-exporter/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Container", func() {
	It("container metrics", func() {
		testCnt01PodName := "exp_e2e_test_cnt01_pod01"
		testCnt01Name := "exp_e2e_test_cnt01"
		testBusyBoxImage := "quay.io/quay/busybox:latest"

		var cnt01Inpect []entities.ContainerInspectReport

		cnt01InspectOutput, err := exec.Command("podman", "container", "inspect", testCnt01Name).Output()
		Expect(err).To(BeNil())
		err = json.Unmarshal(cnt01InspectOutput, &cnt01Inpect)
		Expect(err).To(BeNil())

		cnt01Pod01Inspect, err := utils.PodInformation(testCnt01PodName)
		Expect(err).To(BeNil())

		response := queryEndPoint()

		expectedCnt01Info := fmt.Sprintf("podman_container_info{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		expectedCnt01State := fmt.Sprintf("podman_container_state{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		expectedCnt01Created := fmt.Sprintf("podman_container_created_seconds{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		expectedCnt01ExitedSeconds := fmt.Sprintf("podman_container_exited_seconds{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		expectedCnt01ExitedCode := fmt.Sprintf("podman_container_exit_code{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		expectedCnt01RwSize := fmt.Sprintf("podman_container_rw_size_bytes{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		expectedCnt01RootFsSize := fmt.Sprintf("podman_container_rootfs_size_bytes{id=\"%s\",image=\"%s\",name=\"%s\",pod_id=\"%s\",pod_name=\"%s\",ports=\"\"}",
			cnt01Inpect[0].ID[0:12], testBusyBoxImage, testCnt01Name, cnt01Pod01Inspect.ID, cnt01Pod01Inspect.Name)

		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01Info)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01State)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01Created)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01ExitedSeconds)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01ExitedCode)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01RwSize)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedCnt01RootFsSize)))
	})
})
