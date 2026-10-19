package collector

import (
	"log/slog"

	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/prometheus/client_golang/prometheus"
)

type diskUsageCollector struct {
	imagesSize            typedDesc
	imagesReclaimable     typedDesc
	containersSize        typedDesc
	containersReclaimable typedDesc
	volumesSize           typedDesc
	volumesReclaimable    typedDesc
	logger                *slog.Logger
}

func init() {
	registerCollector("disk_usage", defaultDisabled, NewDiskUsageCollector)
}

// NewDiskUsageCollector returns a Collector exposing podman disk usage information.
func NewDiskUsageCollector(logger *slog.Logger) (Collector, error) {
	return &diskUsageCollector{
		imagesSize: newDiskUsageDesc("images_size_bytes",
			"Podman disk usage of images (sum of all image layers)."),
		imagesReclaimable: newDiskUsageDesc("images_reclaimable_bytes",
			"Podman disk usage estimate of reclaimable image space (images not used by containers)."),
		containersSize: newDiskUsageDesc("containers_size_bytes",
			"Podman disk usage of containers (sum of container read-write layers)."),
		containersReclaimable: newDiskUsageDesc("containers_reclaimable_bytes",
			"Podman disk usage estimate of reclaimable container space (containers not running)."),
		volumesSize: newDiskUsageDesc("volumes_size_bytes",
			"Podman disk usage of local volumes."),
		volumesReclaimable: newDiskUsageDesc("volumes_reclaimable_bytes",
			"Podman disk usage estimate of reclaimable local volume space (volumes not used by containers)."),
		logger: logger,
	}, nil
}

// Update reads and exposes podman disk usage information.
func (c *diskUsageCollector) Update(ch chan<- prometheus.Metric) error {
	diskUsage, err := pdcs.DiskUsageSummary()
	if err != nil {
		return err
	}

	ch <- c.imagesSize.mustNewConstMetric(float64(diskUsage.ImagesSize))

	ch <- c.imagesReclaimable.mustNewConstMetric(float64(diskUsage.ImagesReclaimable))

	ch <- c.containersSize.mustNewConstMetric(float64(diskUsage.ContainersSize))

	ch <- c.containersReclaimable.mustNewConstMetric(float64(diskUsage.ContainersReclaimable))

	ch <- c.volumesSize.mustNewConstMetric(float64(diskUsage.VolumesSize))

	ch <- c.volumesReclaimable.mustNewConstMetric(float64(diskUsage.VolumesReclaimable))

	return nil
}

func newDiskUsageDesc(name string, help string) typedDesc {
	return typedDesc{
		prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "disk_usage", name),
			help,
			nil, nil,
		), prometheus.GaugeValue,
	}
}
