package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type volumeCollector struct {
	info    typedDesc
	created typedDesc
	logger  log.Logger
}

var volumeDefaultLAbels = []string{"name", "driver", "mount_point"}

func init() {
	registerCollector("volume", defaultDisabled, NewVolumeStatsCollector)
}

// NewVolumeStatsCollector returns a Collector exposing volume stats information.
func NewVolumeStatsCollector(logger log.Logger) (Collector, error) {
	createdLabels := []string{"name"}

	if enhanceAllMetrics {
		createdLabels = volumeDefaultLAbels
	}

	return &volumeCollector{
		info: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "volume", "info"),
				"Volume information.",
				volumeDefaultLAbels, nil,
			), prometheus.GaugeValue,
		},
		created: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "volume", "created_seconds"),
				"Volume creation time in unixtime.",
				createdLabels, nil,
			), prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes networks stats.
func (c *volumeCollector) Update(ch chan<- prometheus.Metric) error {
	reports, err := pdcs.Volumes()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		ch <- c.info.mustNewConstMetric(1, rep.Name, rep.Driver, rep.MountPoint)

		if enhanceAllMetrics {
			ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.Name, rep.Driver, rep.MountPoint)

			continue
		}

		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.Name)
	}

	return nil
}
