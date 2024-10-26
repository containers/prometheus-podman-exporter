package pdcs_test

import (
	"os/exec"
	"time"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/promslog"
)

var _ = Describe("Pdcs/Events", func() {
	It("EventStreamer", func() {
		podmanImages, err := pdcs.Images()
		Expect(err).To(BeNil())

		imageCount01 := len(podmanImages)
		logger := promslog.New(&promslog.Config{})
		pdcs.StartEventStreamer(logger, true)

		testImage := "quay.io/libpod/alpine"

		_, err = exec.Command("podman", "image", "pull", testImage).Output()
		Expect(err).To(BeNil())

		time.Sleep(8 * time.Second)

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
