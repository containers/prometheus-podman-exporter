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
			nil, prometheus.GaugeValue,
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
		infoMetric, infoValues := c.getImageInfoDesc(rep)
		c.info.desc = infoMetric
		ch <- c.info.mustNewConstMetric(1, infoValues...)

		ch <- c.size.mustNewConstMetric(float64(rep.Size), rep.ID, rep.Repository, rep.Tag)
		ch <- c.created.mustNewConstMetric(float64(rep.Created), rep.ID, rep.Repository, rep.Tag)
	}

	return nil
}

func (c *imageCollector) getImageInfoDesc(rep pdcs.Image) (*prometheus.Desc, []string) {
	imageLabels := []string{"id", "parent_id", "repository", "tag", "digest"}
	imageLabelsValue := []string{rep.ID, rep.ParentID, rep.Repository, rep.Tag, rep.Digest}

	extraLabels, extraValues := c.getExtraLabelsAndValues(imageLabels, rep)

	imageLabels = append(imageLabels, extraLabels...)
	imageLabelsValue = append(imageLabelsValue, extraValues...)

	infoDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "image", "info"),
		"Image information.",
		imageLabels, nil,
	)

	return infoDesc, imageLabelsValue
}

func (c *imageCollector) getExtraLabelsAndValues(collectorLabels []string, rep pdcs.Image) ([]string, []string) {
	extraLabels := make([]string, 0)
	extraValues := make([]string, 0)

	for label, value := range rep.Labels {
		if slicesContains(collectorLabels, label) {
			continue
		}

		validLabel := sanitizeLabelName(label)
		if storeLabels {
			extraLabels = append(extraLabels, validLabel)
			extraValues = append(extraValues, value)
		} else if whitelistContains(label) {
			extraLabels = append(extraLabels, validLabel)
			extraValues = append(extraValues, value)
		}
	}

	return extraLabels, extraValues
}
