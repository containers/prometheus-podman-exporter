package exporter

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/containers/prometheus-podman-exporter/collector"
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"
)

const minCacheDuration int64 = 5

var errMinCacheDurtion = errors.New("invalid cache duration value, shall be >= " + strconv.Itoa(int(minCacheDuration)))

type exporterOptions struct {
	debug                     bool
	webListen                 string
	webMaxRequests            int
	webTelemetryPath          string
	webDisableExporterMetrics bool
	webConfigFile             string
	cacheDuration             int64
	enableAll                 bool
	storeLabels               bool
	whiteListedLabels         string
	enableImages              bool
	enablePods                bool
	enableVolumes             bool
	enableNetworks            bool
	enableSystem              bool
	enhanceMetrics            bool
}

// Start starts prometheus exporter.
func Start(cmd *cobra.Command, _ []string) error {
	// setup exporter
	logLevel := "info"
	promlogConfig := &promslog.Config{Level: promslog.NewLevel()}

	cmdOptions, err := parseOptions(cmd)
	if err != nil {
		return err
	}

	if cmdOptions.debug {
		logLevel = "debug"
	}

	if err := promlogConfig.Level.Set(logLevel); err != nil {
		return err
	}

	logger := promslog.New(promlogConfig)

	if err := setEnabledCollectors(cmdOptions); err != nil {
		logger.Error("cannot set enabled collectors", "err", err)

		return err
	}

	logger.Info("starting podman-prometheus-exporter", "version", version.Info())
	logger.Info("metrics", "enhanced", cmdOptions.enhanceMetrics)

	http.Handle(
		cmdOptions.webTelemetryPath,
		newHandler(cmdOptions.webDisableExporterMetrics, cmdOptions.webMaxRequests, logger),
	)
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Podman Exporter</title></head>
			<body>
			<h1>Podman Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})

	// setup podman registry
	pdcs.SetupRegistry()
	// start podman event streamer and initiate first update.
	updateImages := false
	if cmdOptions.enableAll || cmdOptions.enableImages {
		updateImages = true
	}

	go pdcs.StartEventStreamer(logger, updateImages)
	pdcs.StartCacheSizeTicker(logger, cmdOptions.cacheDuration)

	logger.Info("Listening on", "address", cmdOptions.webListen)

	server := &http.Server{
		ReadHeaderTimeout: 3 * time.Second, //nolint:mnd
	}
	serverSystemd := false
	serverWebListen := []string{cmdOptions.webListen}

	toolkitFlag := new(web.FlagConfig)
	toolkitFlag.WebSystemdSocket = &serverSystemd
	toolkitFlag.WebListenAddresses = &serverWebListen
	toolkitFlag.WebConfigFile = &cmdOptions.webConfigFile

	if err := web.ListenAndServe(server, toolkitFlag, logger); err != nil {
		return err
	}

	return nil
}

func setEnabledCollectors(opts *exporterOptions) error {
	enList := []string{"container"}

	collector.RegisterVariableLabels(opts.storeLabels, opts.whiteListedLabels, opts.enhanceMetrics)

	if opts.enableAll {
		enList = append(enList, "pod")
		enList = append(enList, "image")
		enList = append(enList, "volume")
		enList = append(enList, "network")
		enList = append(enList, "system")
	} else {
		enList = append(enList, getEnabledCollectors(opts)...)
	}

	// set podman collector state
	for _, col := range enList {
		collector.SetPodmanCollectorState(col, true)
	}

	return nil
}

func getEnabledCollectors(opts *exporterOptions) []string {
	enCollectors := make([]string, 0)

	if opts.enableImages {
		enCollectors = append(enCollectors, "image")
	}

	if opts.enablePods {
		enCollectors = append(enCollectors, "pod")
	}

	if opts.enableVolumes {
		enCollectors = append(enCollectors, "volume")
	}

	if opts.enableNetworks {
		enCollectors = append(enCollectors, "network")
	}

	if opts.enableSystem {
		enCollectors = append(enCollectors, "system")
	}

	return enCollectors
}

func parseOptions(cmd *cobra.Command) (*exporterOptions, error) { //nolint:cyclop
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return nil, err
	}

	webListen, err := cmd.Flags().GetString("web.listen-address")
	if err != nil {
		return nil, err
	}

	webMaxRequests, err := cmd.Flags().GetInt("web.max-requests")
	if err != nil {
		return nil, err
	}

	webTelemetryPath, err := cmd.Flags().GetString("web.telemetry-path")
	if err != nil {
		return nil, err
	}

	webDisableExporterMetrics, err := cmd.Flags().GetBool("web.disable-exporter-metrics")
	if err != nil {
		return nil, err
	}

	webConfigFile, err := cmd.Flags().GetString("web.config.file")
	if err != nil {
		return nil, err
	}

	enableAll, err := cmd.Flags().GetBool("collector.enable-all")
	if err != nil {
		return nil, err
	}

	storeLabels, err := cmd.Flags().GetBool("collector.store_labels")
	if err != nil {
		return nil, err
	}

	whiteListedLabels, err := cmd.Flags().GetString("collector.whitelisted_labels")
	if err != nil {
		return nil, err
	}

	enableImages, err := cmd.Flags().GetBool("collector.image")
	if err != nil {
		return nil, err
	}

	enablePods, err := cmd.Flags().GetBool("collector.pod")
	if err != nil {
		return nil, err
	}

	enableVolumes, err := cmd.Flags().GetBool("collector.volume")
	if err != nil {
		return nil, err
	}

	enableNetworks, err := cmd.Flags().GetBool("collector.network")
	if err != nil {
		return nil, err
	}

	enableSystem, err := cmd.Flags().GetBool("collector.system")
	if err != nil {
		return nil, err
	}

	cacheDuration, err := cmd.Flags().GetInt64("collector.cache_duration")
	if err != nil {
		return nil, err
	}

	if cacheDuration < minCacheDuration {
		return nil, errMinCacheDurtion
	}

	enhanceMetrics, err := cmd.Flags().GetBool("collector.enhance-metrics")
	if err != nil {
		return nil, err
	}

	return &exporterOptions{
		debug:                     debug,
		webListen:                 webListen,
		webMaxRequests:            webMaxRequests,
		webTelemetryPath:          webTelemetryPath,
		webDisableExporterMetrics: webDisableExporterMetrics,
		webConfigFile:             webConfigFile,
		enableAll:                 enableAll,
		storeLabels:               storeLabels,
		whiteListedLabels:         whiteListedLabels,
		enableImages:              enableImages,
		enablePods:                enablePods,
		enableVolumes:             enableVolumes,
		enableNetworks:            enableNetworks,
		enableSystem:              enableSystem,
		cacheDuration:             cacheDuration,
		enhanceMetrics:            enhanceMetrics,
	}, nil
}
