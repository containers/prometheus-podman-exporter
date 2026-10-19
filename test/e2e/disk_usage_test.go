package e2e_test

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type diskUsageSummary struct {
	RawSize        int64
	RawReclaimable int64
}

var diskUsageSystemDf map[string]diskUsageSummary

// createDiskUsageObjects creates podman objects which are counted by disk usage.
// Disk usage is cached on exporter startup, therefore it must be called before the exporter starts.
func createDiskUsageObjects() {
	testBusyBoxImage := "quay.io/quay/busybox"
	testImageTag := "localhost/exp_e2e_test_du_busybox:latest"
	testContainer := "exp_e2e_test_du_cnt01"
	testVolumeUsed := "exp_e2e_test_du_vol01"
	testVolumeUnused := "exp_e2e_test_du_vol02"

	_, err := exec.Command("podman", "image", "pull", testBusyBoxImage).Output()
	Expect(err).To(BeNil())

	// second tag for the same image, image shall be counted only once.
	_, err = exec.Command("podman", "image", "tag", testBusyBoxImage, testImageTag).Output()
	Expect(err).To(BeNil())

	_, err = exec.Command("podman", "volume", "create", testVolumeUnused).Output()
	Expect(err).To(BeNil())

	// container uses the image and a volume.
	_, err = exec.Command("podman", "container", "create", "--name", testContainer,
		"-v", testVolumeUsed+":/data", testBusyBoxImage).Output()
	Expect(err).To(BeNil())
}

func podmanSystemDf() map[string]diskUsageSummary {
	output, err := exec.Command("podman", "system", "df", "--format", "json").Output()
	Expect(err).To(BeNil())

	var reports []struct {
		Type string
		diskUsageSummary
	}

	Expect(json.Unmarshal(output, &reports)).To(Succeed())

	summaries := make(map[string]diskUsageSummary)
	for _, report := range reports {
		summaries[report.Type] = report.diskUsageSummary
	}

	return summaries
}

func metricValue(response []string, metric string) float64 {
	for _, line := range response {
		value, found := strings.CutPrefix(line, metric+" ")
		if !found {
			continue
		}

		parsed, err := strconv.ParseFloat(value, 64)
		Expect(err).To(BeNil())

		return parsed
	}

	Fail("metric not found: " + metric)

	return 0
}

var _ = Describe("DiskUsage", func() {
	It("disk usage metrics", func() {
		response := queryEndPoint()

		Expect(diskUsageSystemDf).To(HaveKey("Images"))
		Expect(diskUsageSystemDf).To(HaveKey("Containers"))
		Expect(diskUsageSystemDf).To(HaveKey("Local Volumes"))
		Expect(diskUsageSystemDf["Images"].RawSize).To(BeNumerically(">", 0))

		expectedValues := map[string]diskUsageSummary{
			"images":     diskUsageSystemDf["Images"],
			"containers": diskUsageSystemDf["Containers"],
			"volumes":    diskUsageSystemDf["Local Volumes"],
		}

		for name, expected := range expectedValues {
			Expect(metricValue(response, "podman_disk_usage_"+name+"_size_bytes")).
				To(Equal(float64(expected.RawSize)), name+" size")
			Expect(metricValue(response, "podman_disk_usage_"+name+"_reclaimable_bytes")).
				To(Equal(float64(expected.RawReclaimable)), name+" reclaimable")
		}
	})
})
