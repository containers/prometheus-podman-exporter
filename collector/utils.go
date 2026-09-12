package collector

import (
	"regexp"
	"slices"
	"strings"
	"sync"
)

const (
	labelName       = "name"
	labelID         = "id"
	labelImage      = "image"
	labelImageID    = "image_id"
	labelPorts      = "ports"
	labelPodID      = "pod_id"
	labelPodName    = "pod_name"
	labelDriver     = "driver"
	labelInterface  = "interface"
	labelLabels     = "labels"
	labelInfraID    = "infra_id"
	labelParentID   = "parent_id"
	labelRepository = "repository"
	labelTag        = "tag"
	labelDigest     = "digest"
	labelMountPoint = "mount_point"
	labelVersion    = "version"
)

var (
	collectorSync          sync.Once
	storeLabels            bool
	enhanceAllMetrics      bool
	whitelistedLabels      []string
	invalidNameCharRE      = regexp.MustCompile(`[^a-zA-Z0-9_]`)
	containerDefaultLabels = []string{labelID, labelName, labelImage, labelImageID, labelPorts, labelPodID, labelPodName}
	podDefaultLabels       = []string{labelID}
	podDescDefaultLabels   = []string{labelID, labelName, labelInfraID}
	networkDefaultLabels   = []string{labelName, labelID, labelDriver, labelInterface, labelLabels}
	imageDefaultLabels     = []string{labelID, labelRepository, labelTag}
	imageDescDefaultLabels = []string{labelID, labelParentID, labelRepository, labelTag, labelDigest}
	volumeDefaultLabels    = []string{labelName, labelDriver, labelMountPoint}
	systemDefaultLabels    = []string{labelVersion}
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
	return slices.Contains(whitelistedLabels, text)
}

func slicesContains(list []string, value string) bool {
	val := strings.ToLower(value)

	return slices.Contains(list, val)
}
