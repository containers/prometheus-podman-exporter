package pdcs_test

import (
	"encoding/json"
	"os/exec"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

var _ = Describe("Pdcs/DiskUsage", func() {
	It("SummarizeDiskUsage", func() {
		report := &entities.SystemDfReport{
			ImagesSize: 200,
			Images: []*entities.SystemDfImageReport{
				{ImageID: "used", Containers: 1, UniqueSize: 40},
				{ImageID: "used", Containers: 1, UniqueSize: 40}, // second tag of the same image
				{ImageID: "unused", Containers: 0, UniqueSize: 25},
			},
			Containers: []*entities.SystemDfContainerReport{
				{Status: "running", RWSize: 10},
				{Status: "created", RWSize: 5},
				{Status: "exited", RWSize: 7},
			},
			Volumes: []*entities.SystemDfVolumeReport{
				{Size: 30, ReclaimableSize: 0},
				{Size: 50, ReclaimableSize: 50},
			},
		}

		Expect(pdcs.SummarizeDiskUsage(report)).To(Equal(pdcs.DiskUsage{
			ImagesSize:            200,
			ImagesReclaimable:     160, // used image subtracted once
			ContainersSize:        22,
			ContainersReclaimable: 12, // running container excluded
			VolumesSize:           80,
			VolumesReclaimable:    50,
		}))
	})

	It("DiskUsageSummary", func() {
		testImage := "quay.io/quay/busybox"
		testImageTag := "localhost/exp_pdcs_test_du_busybox:latest"
		testContainer := "exp_pdcs_test_du_container01"
		testVolumeUsed := "exp_pdcs_test_du_vol01"
		testVolumeUnused := "exp_pdcs_test_du_vol02"

		DeferCleanup(func() {
			_ = exec.Command("podman", "container", "rm", "-f", "-t", "0", testContainer).Run()
			_ = exec.Command("podman", "volume", "rm", "-f", testVolumeUsed, testVolumeUnused).Run()
			_ = exec.Command("podman", "image", "untag", testImage, testImageTag).Run()
		})

		_, err := exec.Command("podman", "image", "pull", testImage).Output()
		Expect(err).To(BeNil())

		// second tag for the same image, image shall be counted only once.
		_, err = exec.Command("podman", "image", "tag", testImage, testImageTag).Output()
		Expect(err).To(BeNil())

		_, err = exec.Command("podman", "volume", "create", testVolumeUnused).Output()
		Expect(err).To(BeNil())

		// container uses the image and a volume.
		_, err = exec.Command("podman", "container", "create", "--name", testContainer,
			"-v", testVolumeUsed+":/data", testImage).Output()
		Expect(err).To(BeNil())

		pdcs.UpdateDiskUsage()

		diskUsage, err := pdcs.DiskUsageSummary()
		Expect(err).To(BeNil())

		output, err := exec.Command("podman", "system", "df", "--format", "json").Output()
		Expect(err).To(BeNil())

		type summary struct {
			RawSize        int64
			RawReclaimable int64
		}

		var reports []struct {
			Type string
			summary
		}

		Expect(json.Unmarshal(output, &reports)).To(Succeed())

		summaries := make(map[string]summary)
		for _, report := range reports {
			summaries[report.Type] = report.summary
		}

		Expect(summaries).To(HaveKey("Images"))
		Expect(summaries).To(HaveKey("Containers"))
		Expect(summaries).To(HaveKey("Local Volumes"))

		Expect(diskUsage.ImagesSize).To(BeNumerically(">", 0))

		Expect(summary{diskUsage.ImagesSize, diskUsage.ImagesReclaimable}).
			To(Equal(summaries["Images"]))
		Expect(summary{diskUsage.ContainersSize, diskUsage.ContainersReclaimable}).
			To(Equal(summaries["Containers"]))
		Expect(summary{diskUsage.VolumesSize, diskUsage.VolumesReclaimable}).
			To(Equal(summaries["Local Volumes"]))
	})
})
