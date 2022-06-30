package pdcs

import (
	"strings"

	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/libpod/define"
)

// System implements podman system information.
type System struct {
	Podman  string
	Runtime string
	Conmon  string
	Buildah string
}

type reportInfo struct {
	*define.Info
}

// SystemInfo return system information (System).
func SystemInfo() (System, error) {
	var sysinfo System

	report, err := registry.ContainerEngine().Info(registry.GetContext())
	if err != nil {
		return sysinfo, err
	}

	sysinfo.Podman = report.Version.APIVersion
	sysinfo.Runtime = reportInfo{report}.runtimeVersion()
	sysinfo.Conmon = reportInfo{report}.conmonVersion()
	sysinfo.Buildah = report.Host.BuildahVersion

	return sysinfo, nil
}

func (r reportInfo) runtimeVersion() string {
	runtime := strings.Split(r.Host.OCIRuntime.Version, ":")[0]
	runtime = strings.ReplaceAll(runtime, "commit", "")
	runtime = strings.Trim(runtime, "\n")

	return runtime
}

func (r reportInfo) conmonVersion() string {
	conmonVersion := strings.Split(r.Host.Conmon.Version, ",")[0]
	conmonVersion = strings.ReplaceAll(conmonVersion, "conmon version", "")
	conmonVersion = strings.TrimSpace(conmonVersion)

	return conmonVersion
}
