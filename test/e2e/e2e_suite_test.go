package e2e_test

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/containers/prometheus-podman-exporter/exporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var (
	endpointURL         = "http://127.0.0.1:9882/metrics"
	cacheDuration int64 = 3600
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var _ = BeforeSuite(func() {
	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "",
		Long:  ``,
	}

	rootCmd.Flags().BoolP("debug", "d", true, "")
	rootCmd.Flags().BoolP("version", "", false, "")
	rootCmd.Flags().StringP("web.config.file", "", "", "")
	rootCmd.Flags().StringArrayP("web.listen-addresses", "l", []string{":9882"}, "")
	rootCmd.Flags().StringP("web.telemetry-path", "p", "/metrics", "")
	rootCmd.Flags().BoolP("web.disable-exporter-metrics", "e", false, "")
	rootCmd.Flags().IntP("web.max-requests", "m", 10, "")
	rootCmd.Flags().BoolP("collector.enable-all", "a", true, "")
	rootCmd.Flags().BoolP("collector.image", "i", false, "")
	rootCmd.Flags().BoolP("collector.pod", "o", false, "")
	rootCmd.Flags().BoolP("collector.volume", "v", false, "")
	rootCmd.Flags().BoolP("collector.network", "n", false, "")
	rootCmd.Flags().BoolP("collector.system", "s", false, "")
	rootCmd.Flags().BoolP("collector.store_labels", "b", false, "")
	rootCmd.Flags().StringP("collector.whitelisted_labels", "w", "", "")
	rootCmd.Flags().Int64P("collector.cache_duration", "t", cacheDuration, "")
	rootCmd.Flags().BoolP("collector.enhance-metrics", "", false, "")

	go func() {
		err := exporter.Start(rootCmd, nil)
		Expect(err).To(BeNil())
	}()

	time.Sleep(10 * time.Second)
})

func extractLabelValue(line string, label string) string {
	value := "-9999999999999"
	r := regexp.MustCompile("{(.*)}")
	matches := r.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return value
	}

	if len(matches[0]) == 0 {
		return value
	}

	for _, item := range strings.Split(matches[0][1], ",") {
		if strings.Index(item, fmt.Sprintf("%s=", label)) == 0 {
			return strings.ReplaceAll(strings.Split(item, "=")[1], "\"", "")
		}
	}

	return value
}

func queryEndPoint() []string {
	req, err := http.NewRequest("GET", endpointURL, nil)
	Expect(err).To(BeNil())

	res, err := http.DefaultClient.Do(req)
	Expect(err).To(BeNil())

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	Expect(err).To(BeNil())

	return strings.Split(string(body), "\n")
}
