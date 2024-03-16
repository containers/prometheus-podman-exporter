package collector

import (
	"regexp"
	"strings"
	"sync"
)

var (
	collectorSync     sync.Once
	storeLabels       bool
	enhanceAllMetrics bool
	whitelistedLabels []string
	invalidNameCharRE = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

// RegisterVariableLabels sets storeLabels or whiteListed labels to be converted to metrics.
func RegisterVariableLabels(storeLabel bool, whiteListed string, enhanceMetrics bool) {
	collectorSync.Do(func() {
		storeLabels = storeLabel
		whitelistedLabels = strings.Split(whiteListed, ",")
		enhanceAllMetrics = enhanceMetrics
	})
}

func sanitizeLabelName(name string) string {
	return invalidNameCharRE.ReplaceAllString(name, "_")
}

func whitelistContains(text string) bool {
	for _, item := range whitelistedLabels {
		if item == text {
			return true
		}
	}

	return false
}

func slicesContains(list []string, value string) bool {
	for _, item := range list {
		if item == strings.ToLower(value) {
			return true
		}
	}

	return false
}
