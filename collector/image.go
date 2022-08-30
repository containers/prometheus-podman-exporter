package collector

import (
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type imageCollector struct {
	info    typedDesc
	size    typedDesc
	created typedDesc
	logger  log.Logger
}

func init() {
	registerCollector("image", defaultDisabled, NewImageStatsCollector)
}

// NewImageStatsCollector returns a Collector exposing image stats information.
func NewImageStatsCollector(logger log.Logger) (Collector, error) {
	return &imageCollector{
		info: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "image", "info"),
				"Image information.",
				[]string{"id", "parent_id", "repository", "tag"}, nil,
			), prometheus.GaugeValue,
		},
		size: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "image", "size"),
				"Image size",
				[]string{"id", "repository", "tag"}, nil,
			), prometheus.GaugeValue,
		},
		created: typedDesc{
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "image", "created_seconds"),
				"Image creation time in unixtime.",
				[]string{"id", "repository", "tag"}, nil,
			), prometheus.GaugeValue,
		},
		logger: logger,
	}, nil
}

// Update reads and exposes images stats.
func (c *imageCollector) Update(ch chan<- prometheus.Metric) error {
	reports, err := pdcs.Images()
	if err != nil {
		return err
	}

	for _, rep := range reports {
		ch <- c.info.mustNewConstMetric(1, rep.ID, rep.ParentID, rep.Repository, rep.Tag)
		ch <- c.size.mustNewConstMetric(float64(rep.Size), rep.ID, rep.Repository, rep.Tag)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID, rep.Repository, rep.Tag)
	}

	return nil
}
