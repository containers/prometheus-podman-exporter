package e2e_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("System", func() {
	It("system metrics", func() {
		response := queryEndPoint()

		apiVersion := ""
		buildahVersion := ""
		conmonVersion := ""
		runtimeVersion := ""

		for _, line := range response {
			if strings.Index(line, "podman_system_api_version") == 0 {
				apiVersion = extractLabelValue(line, "version")

				continue
			}

			if strings.Index(line, "podman_system_buildah_version") == 0 {
				buildahVersion = extractLabelValue(line, "version")

				continue
			}

			if strings.Index(line, "podman_system_conmon_version") == 0 {
				conmonVersion = extractLabelValue(line, "version")

				continue
			}

			if strings.Index(line, "podman_system_runtime_version") == 0 {
				runtimeVersion = extractLabelValue(line, "version")

				continue
			}
		}

		Expect(apiVersion).NotTo(BeEmpty())
		Expect(buildahVersion).NotTo(BeEmpty())
		Expect(conmonVersion).NotTo(BeEmpty())
		Expect(runtimeVersion).NotTo(BeEmpty())
	})
})
