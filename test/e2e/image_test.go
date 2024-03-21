package e2e_test

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/containers/podman/v5/pkg/domain/entities"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image", func() {
	It("image metrics", func() {
		testBusyBoxImage := "quay.io/quay/busybox"

		_, err := exec.Command("podman", "image", "pull", testBusyBoxImage).Output()
		Expect(err).To(BeNil())

		imageInpsectOutput, err := exec.Command("podman", "image", "inspect", testBusyBoxImage).Output()
		Expect(err).To(BeNil())

		var imageInspect []entities.ImageInspectReport
		err = json.Unmarshal(imageInpsectOutput, &imageInspect)
		Expect(err).To(BeNil())

		response := queryEndPoint()
		expectedImageSize := fmt.Sprintf("podman_image_size{id=\"%s\",repository=\"%s\",tag=\"latest\"}",
			imageInspect[0].ID[0:12], testBusyBoxImage)

		expectedImageCreated := fmt.Sprintf("podman_image_created_seconds{id=\"%s\",repository=\"%s\",tag=\"latest\"}",
			imageInspect[0].ID[0:12], testBusyBoxImage)

		expectedImageInfo := fmt.Sprintf("podman_image_info{digest=\"%s\",id=\"%s\",parent_id=\"\",repository=\"%s\",tag=\"latest\"} 1",
			imageInspect[0].Digest.String(), imageInspect[0].ID[0:12], testBusyBoxImage)

		Expect(response).Should(ContainElement(ContainSubstring(expectedImageSize)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedImageCreated)))
		Expect(response).Should(ContainElement(ContainSubstring(expectedImageInfo)))
	})
})
