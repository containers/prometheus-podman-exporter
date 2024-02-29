package pdcs_test

import (
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/promlog"
)

var _ = Describe("Pdcs/Events", func() {
	It("EventStreamer", func() {
		podmanImages, err := pdcs.Images()
		Expect(err).To(BeNil())

		imageCount01 := len(podmanImages)
		logger := promlog.New(&promlog.Config{})
		pdcs.StartEventStreamer(logger)

		testImage := "docker.io/library/alpine"

		_, err = exec.Command("podman", "image", "pull", testImage).Output()
		Expect(err).To(BeNil())

		podmanImages, err = pdcs.Images()
		Expect(err).To(BeNil())

		imageCount02 := len(podmanImages)

		imageIsUpdated := false
		if imageCount02 > imageCount01 {
			imageIsUpdated = true
		}

		Expect(imageIsUpdated).To(BeTrue())
	})
})
