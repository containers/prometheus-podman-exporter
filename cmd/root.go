package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/containers/prometheus-podman-exporter/exporter"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
)

const (
	maxRequest int = 40
)

var (
	buildVersion  string
	buildRevision string
	buildBranch   string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "prometheus-podman-exporter",
	Short: "prometheus-podman-exporter",
	Long: `Prometheus exporter for podman exposing containers, pods, images,
volumes and networks information.`,
	PreRunE: preRun,
	Run:     run,
}

func preRun(cmd *cobra.Command, args []string) error {
	version.Version = buildVersion
	version.Revision = buildRevision
	version.Branch = buildBranch

	printVersion, err := cmd.Flags().GetBool("version")
	if err != nil {
		return err
	}

	if printVersion {
		fmt.Println(cmd.Use, version.Info())
		os.Exit(1)
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	if err := exporter.Start(cmd, args); err != nil {
		log.Panic(err.Error())
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("debug", "d", false,
		"Set log level to debug.")
	rootCmd.Flags().BoolP("version", "", false,
		"Print version and exit.")
	rootCmd.Flags().StringP("web.listen-address", "l", ":9882",
		"Address on which to expose metrics and web interface.")
	rootCmd.Flags().StringP("web.telemetry-path", "p", "/metrics",
		"Path under which to expose metrics.")
	rootCmd.Flags().BoolP("web.disable-exporter-metrics", "e", false,
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).")
	rootCmd.Flags().IntP("web.max-requests", "m", maxRequest,
		"Maximum number of parallel scrape requests. Use 0 to disable")
	rootCmd.Flags().BoolP("collector.enable-all", "a", false,
		"Enable all collectors by default.")
	rootCmd.Flags().BoolP("collector.image", "i", false,
		"Enable image collector.")
	rootCmd.Flags().BoolP("collector.pod", "o", false,
		"Enable pod collector.")
	rootCmd.Flags().BoolP("collector.volume", "v", false,
		"Enable volume collector.")
	rootCmd.Flags().BoolP("collector.network", "n", false,
		"Enable network collector.")
	rootCmd.Flags().BoolP("collector.system", "s", false,
		"Enable system collector.")
}
