package pdcs_test

import (
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pdcs/Image", func() {
	It("Images", func() {
		testImage := "docker.io/library/busybox"

		exec.Command("podman", "image", "rm", testImage)
		_, err := exec.Command("podman", "image", "pull", testImage).Output()
		Expect(err).To(BeNil())

		pdcs.UpdateImages()

		podmanImages, err := pdcs.Images()
		Expect(err).To(BeNil())

		imageFound := false
		for _, image := range podmanImages {
			if image.Repository == testImage {
				imageFound = true

				break
			}
		}

		Expect(imageFound).To(BeTrue())
	})
})
