package exporter

import (
	"net/http"

	"github.com/containers/prometheus-podman-exporter/collector"
	"github.com/containers/prometheus-podman-exporter/pdcs"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"
)

// Start starts prometheus exporter.
func Start(cmd *cobra.Command, args []string) error {
	// setup podman resgistry
	pdcs.SetupRegistry()

	// setup exporter
	promlogConfig := &promlog.Config{Level: &promlog.AllowedLevel{}}

	logLevel := "info"

	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		logLevel = "debug"
	}

	if err := promlogConfig.Level.Set(logLevel); err != nil {
		return err
	}

	webListen, err := cmd.Flags().GetString("web.listen-address")
	if err != nil {
		return err
	}

	webMaxRequests, err := cmd.Flags().GetInt("web.max-requests")
	if err != nil {
		return err
	}

	webTelemetryPath, err := cmd.Flags().GetString("web.telemetry-path")
	if err != nil {
		return err
	}

	webDisableExporterMetrics, err := cmd.Flags().GetBool("web.disable-exporter-metrics")
	if err != nil {
		return err
	}

	logger := promlog.New(promlogConfig)

	if err := setEnabledCollectors(cmd); err != nil {
		level.Error(logger).Log("msg", "cannot set enabled collectors", "err", err)

		return err
	}

	level.Info(logger).Log("msg", "Starting podman-prometheus-exporter", "version", version.Info())
	http.Handle(webTelemetryPath, newHandler(webDisableExporterMetrics, webMaxRequests, logger))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Podman Exporter</title></head>
			<body>
			<h1>Podman Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	level.Info(logger).Log("msg", "Listening on", "address", webListen)

	server := &http.Server{}
	serverSystemd := false
	serverConfigFile := ""
	serverWebListen := []string{webListen}

	toolkitFlag := new(web.FlagConfig)
	toolkitFlag.WebSystemdSocket = &serverSystemd
	toolkitFlag.WebListenAddresses = &serverWebListen
	toolkitFlag.WebConfigFile = &serverConfigFile

	if err := web.ListenAndServe(server, toolkitFlag, logger); err != nil {
		return err
	}

	return nil
}

func setEnabledCollectors(cmd *cobra.Command) error {
	enList := []string{"container"}

	enableAll, err := cmd.Flags().GetBool("collector.enable-all")
	if err != nil {
		return err
	}

	if enableAll {
		enList = append(enList, "pod")
		enList = append(enList, "image")
		enList = append(enList, "volume")
		enList = append(enList, "network")
		enList = append(enList, "system")
	} else {
		enList = append(enList, getEnabledCollectors(cmd)...)
	}

	// set podman collector state
	for _, col := range enList {
		collector.SetPodmanCollectorState(col, true)
	}

	return nil
}

func getEnabledCollectors(cmd *cobra.Command) []string {
	enCollectors := make([]string, 0)

	enimage := command{cmd}.isEnabled("collector.image")
	if enimage {
		enCollectors = append(enCollectors, "image")
	}

	enpod := command{cmd}.isEnabled("collector.pod")
	if enpod {
		enCollectors = append(enCollectors, "pod")
	}

	envolume := command{cmd}.isEnabled("collector.volume")
	if envolume {
		enCollectors = append(enCollectors, "volume")
	}

	ennetwork := command{cmd}.isEnabled("collector.network")
	if ennetwork {
		enCollectors = append(enCollectors, "network")
	}

	ensystem := command{cmd}.isEnabled("collector.system")
	if ensystem {
		enCollectors = append(enCollectors, "system")
	}

	return enCollectors
}

type command struct {
	*cobra.Command
}

func (c command) isEnabled(name string) bool {
	enable, err := c.Flags().GetBool(name)
	if err != nil {
		return false
	}

	return enable
}
