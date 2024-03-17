package e2e_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/podman/v5/pkg/domain/entities"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Volume", func() {
	It("volume metrics", func() {
		testVolName := "exp_e2e_test_vol01"

		_, err := exec.Command("podman", "volume", "create", testVolName).Output()
		Expect(err).To(BeNil())

		volInspectOutput, err := exec.Command("podman", "volume", "inspect", testVolName).Output()
		Expect(err).To(BeNil())

		var volInspect []entities.VolumeInspectReport

		err = json.Unmarshal(volInspectOutput, &volInspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()
		expectedVolCreated := fmt.Sprintf("podman_volume_created_seconds{name=\"%s\"}", testVolName)
		expectedVolInfo := fmt.Sprintf("podman_volume_info{driver=\"%s\",mount_point=\"%s\",name=\"%s\"} 1", volInspect[0].Driver, volInspect[0].Mountpoint, testVolName)

		Expect(response).Should(ContainElement(ContainSubstring(expectedVolCreated)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedVolInfo)))
	})
})
