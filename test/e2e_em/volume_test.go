package e2e_em_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/podman/v4/pkg/domain/entities"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Volume", func() {
	It("volume metrics", func() {
		testVolName := "exp_e2e_test_vol01"

		volInspectOutput, err := exec.Command("podman", "volume", "inspect", testVolName).Output()
		Expect(err).To(BeNil())

		var volInspect []entities.VolumeInspectReport

		err = json.Unmarshal(volInspectOutput, &volInspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()
		expectedVolCreated := fmt.Sprintf("podman_volume_created_seconds{driver=\"%s\",mount_point=\"%s\",name=\"%s\"}", volInspect[0].Driver, volInspect[0].Mountpoint, testVolName)
		expectedVolInfo := fmt.Sprintf("podman_volume_info{driver=\"%s\",mount_point=\"%s\",name=\"%s\"}", volInspect[0].Driver, volInspect[0].Mountpoint, testVolName)

		Expect(response).Should(ContainElement(ContainSubstring(expectedVolCreated)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedVolInfo)))
	})
})
