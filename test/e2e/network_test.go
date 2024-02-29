package e2e_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/common/libnetwork/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Network", func() {
	It("network metrics", func() {
		testNetName := "exp_e2e_test_net01"

		_, err := exec.Command("podman", "network", "create", testNetName).Output()
		Expect(err).To(BeNil())

		netInspectOutput, err := exec.Command("podman", "network", "inspect", testNetName).Output()
		Expect(err).To(BeNil())

		var netInspect []types.Network

		err = json.Unmarshal(netInspectOutput, &netInspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()
		expectedNetworkInfo := fmt.Sprintf("podman_network_info{driver=\"%s\",id=\"%s\",interface=\"%s\",labels=\"\",name=\"%s\"} 1",
			netInspect[0].Driver, netInspect[0].ID[0:12], netInspect[0].NetworkInterface, testNetName)

		Expect(response).Should(ContainElement(ContainSubstring(expectedNetworkInfo)))
	})
})
