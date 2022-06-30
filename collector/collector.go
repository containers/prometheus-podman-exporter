package collector

import (
	"errors"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// Namespace defines the common namespace to be used by all metrics.
const namespace = "podman"

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"),
		"podman_prometheus_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "scrape", "collector_success"),
		"podman_prometheus_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)

	// ErrNoData indicates the collector found no data to collect, but had no other error.
	ErrNoData = errors.New("collector returned no data")
)

const (
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	factories              = make(map[string]func(logger log.Logger) (Collector, error))
	initiatedCollectorsMtx = sync.Mutex{}
	initiatedCollectors    = make(map[string]Collector)
	collectorState         = make(map[string]bool)
)

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(ch chan<- prometheus.Metric) error
}

type typedDesc struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
}

// PodmanCollector implements the prometheus.Collector interface.
type PodmanCollector struct {
	Collectors map[string]Collector
	logger     log.Logger
}

// NewPodmanCollector creates a new PodmanCollector.
func NewPodmanCollector(logger log.Logger) (*PodmanCollector, error) {
	collectors := make(map[string]Collector)

	initiatedCollectorsMtx.Lock()
	defer initiatedCollectorsMtx.Unlock()

	for key, enabled := range collectorState {
		if !enabled {
			continue
		}

		if collector, ok := initiatedCollectors[key]; ok {
			collectors[key] = collector
		} else {
			collector, err := factories[key](log.With(logger, "collector", key))
			if err != nil {
				return nil, err
			}
			collectors[key] = collector
			initiatedCollectors[key] = collector
		}
	}

	return &PodmanCollector{Collectors: collectors, logger: logger}, nil
}

// Describe implements the prometheus.Collector interface.
func (p PodmanCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (p PodmanCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(p.Collectors))

	for name, c := range p.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch, p.logger)
			wg.Done()
		}(name, c)
	}

	wg.Wait()
}

func execute(name string, c Collector, ch chan<- prometheus.Metric, logger log.Logger) {
	var success float64

	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)

	if err != nil {
		if IsNoDataError(err) {
			level.Debug(logger).Log("msg", "collector returned no data", "name",
				name,
				"duration_seconds",
				duration.Seconds(),
				"err", err)
		} else {
			level.Error(logger).Log("msg", "collector failed", "name",
				name,
				"duration_seconds",
				duration.Seconds(),
				"err", err)
		}

		success = 0
	} else {
		level.Debug(logger).Log("msg", "collector succeeded", "name",
			name,
			"duration_seconds",
			duration.Seconds())
		success = 1
	}

	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

// IsNoDataError returns true if error is no data error.
func IsNoDataError(err error) bool {
	return errors.Is(err, ErrNoData)
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

func registerCollector(collector string, enabledByDefault bool, factory func(logger log.Logger) (Collector, error)) {
	collectorState[collector] = enabledByDefault
	factories[collector] = factory
}

// SetPodmanCollectorState enable/disable collectors.
func SetPodmanCollectorState(name string, state bool) {
	for key := range collectorState {
		if key == name {
			collectorState[key] = state
		}
	}
}
