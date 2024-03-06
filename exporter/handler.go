package exporter

import (
	"fmt"
	stdlog "log"
	"net/http"
	"sort"

	"github.com/containers/prometheus-podman-exporter/collector"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

// handler wraps an unfiltered http.Handler but uses a filtered handler,
// created on the fly, if filtering is requested. Create instances with
// newHandler.
type handler struct {
	unfilteredHandler http.Handler
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	disableExporterMetrics  bool
	exporterMetricsRegistry *prometheus.Registry
	maxRequests             int
	logger                  log.Logger
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	level.Debug(h.logger).Log("msg", "collect query:", "filters", fmt.Sprintf("%v", filters))

	if len(filters) == 0 {
		// No filters, use the prepared unfiltered handler.
		h.unfilteredHandler.ServeHTTP(w, r)

		return
	}
	// To serve filtered metrics, we create a filtering handler on the fly.
	filteredHandler, err := h.innerHandler(filters...)
	if err != nil {
		level.Warn(h.logger).Log("msg", "couldn't create filtered metrics handler:", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("couldn't create filtered metrics handler: %s", err)))

		return
	}

	filteredHandler.ServeHTTP(w, r)
}

func newHandler(disableExporterMetrics bool, maxRequests int, logger log.Logger) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		maxRequests:             maxRequests,
		disableExporterMetrics:  disableExporterMetrics,
		logger:                  logger,
	}

	if !disableExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}

	innerHandler, err := h.innerHandler()
	if err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	}

	h.unfilteredHandler = innerHandler

	return h
}

// innerHandler is used to create both the one unfiltered http.Handler to be
// wrapped by the outer handler and also the filtered handlers created on the
// fly. The former is accomplished by calling innerHandler without any arguments
// (in which case it will log all the collectors enabled via command-line
// flags).
func (h *handler) innerHandler(filters ...string) (http.Handler, error) {
	podc, err := collector.NewPodmanCollector(h.logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %w", err)
	}

	// Only log the creation of an unfiltered handler, which should happen
	// only once upon startup.
	if len(filters) == 0 {
		level.Info(h.logger).Log("msg", "enabled collectors")

		collectors := []string{}

		for n := range podc.Collectors {
			collectors = append(collectors, n)
		}

		sort.Strings(collectors)

		for _, c := range collectors {
			level.Info(h.logger).Log("collector", c)
		}
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("prometheus_podman_exporter"))

	if err := r.Register(podc); err != nil {
		return nil, fmt.Errorf("couldn't register podman collector: %w", err)
	}

	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(h.logger)), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: h.maxRequests,
			Registry:            h.exporterMetricsRegistry,
		},
	)

	if !h.disableExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	}

	return handler, nil
}
