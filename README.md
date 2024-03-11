# prometheus-podman-exporter

[![PkgGoDev](https://pkg.go.dev/badge/github.com/containers/prometheus-podman-exporter)](https://pkg.go.dev/github.com/containers/prometheus-podman-exporter)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](https://goreportcard.com/report/github.com/containers/prometheus-podman-exporter)
![Go](https://github.com/containers/prometheus-podman-exporter/workflows/Go/badge.svg)

Prometheus exporter for podman v4.x environment exposing containers, pods, images, volumes and networks information.

prometheus-podman-exporter uses the podman v4.x (libpod) library to fetch the statistics and therefore no need to enable podman.socket service unless using the container image.

- [**Installation**](#installation)
- [**Usage and Options**](#usage-and-options)
- [**Collectors**](#collectors)
- [**License**](#license)

## Installation

Building from source, using container image or installing packaged versions are detailed in [install guide](install.md).

## Usage and Options

```shell
Usage:
  prometheus-podman-exporter [flags]

Flags:
  -t, --collector.cache_duration int          Duration (seconds) to retrieve container, size and refresh the cache (default 3600)
  -a, --collector.enable-all                  Enable all collectors by default.
  -i, --collector.image                       Enable image collector.
  -n, --collector.network                     Enable network collector.
  -o, --collector.pod                         Enable pod collector.
  -b, --collector.store_labels                Convert pod/container/image labels on prometheus metrics for each pod/container/image.
  -s, --collector.system                      Enable system collector.
  -v, --collector.volume                      Enable volume collector.
  -w, --collector.whitelisted_labels string   Comma separated list of pod/container/image labels to be converted
                                              to labels on prometheus metrics for each pod/container/image.
                                              collector.store_labels must be set to false for this to take effect.
  -d, --debug                                 Set log level to debug.
  -h, --help                                  help for prometheus-podman-exporter
      --version                               Print version and exit.
      --web.config.file string                [EXPERIMENTAL] Path to configuration file that can enable TLS or authentication.
  -e, --web.disable-exporter-metrics          Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).
  -l, --web.listen-address string             Address on which to expose metrics and web interface. (default ":9882")
  -m, --web.max-requests int                  Maximum number of parallel scrape requests. Use 0 to disable (default 40)
  -p, --web.telemetry-path string             Path under which to expose metrics. (default "/metrics")
```

By default only container collector is enabled, in order to enable all collectors use `--collector.enable-all` or use `--collector.enable-<name>` flag to enable other collector.

`Example:` enable all available collectors:

```shell
$ ./bin/prometheus-podman-exporter --collector.enable-all
```

The exporter uses plain HTTP without any form of authentication to expose the metrics by default.
Use `--web.config.file` with a configuration file to use TLS for confidentiality and/or to enable authentication.
Visit [this page](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md) for more information about the syntax of the configuration file.

## Collectors
The table below list all existing collector and their description.

| Name      | Description |
| --------- | ----------- |
| container | exposes containers information
| image     | exposes images information
| network   | exposes networks information
| pod       | exposes pod information
| volume    | exposes volume information
| system    | exposes system (host) information

### Collectors examples output

#### `container`

```shell
# HELP podman_container_info Container information.
# TYPE podman_container_info gauge
podman_container_info{id="19286a13dc23",image="docker.io/library/sonarqube:latest",name="sonar01",pod_id="",pod_name="",ports="0.0.0.0:9000->9000/tcp"} 1
podman_container_info{id="482113b805f7",image="docker.io/library/httpd:latest",name="web_server",pod_id="",pod_name="",ports="0.0.0.0:8000->80/tcp"} 1
podman_container_info{id="642490688d9c",image="docker.io/grafana/grafana:latest",name="grafana",pod_id="",pod_name="",ports="0.0.0.0:3000->3000/tcp"} 1
podman_container_info{id="ad36e85960a1",image="docker.io/library/busybox:latest",name="busybox01",pod_id="3e8bae64e9af",pod_name="pod01",ports=""} 1
podman_container_info{id="dda983cc3ecf",image="localhost/podman-pause:4.1.0-1651853754",name="3e8bae64e9af-infra",pod_id="3e8bae64e9af",pod_name="pod01",ports=""} 1

# HELP podman_container_state Container current state (-1=unknown,0=created,1=initialized,2=running,3=stopped,4=paused,5=exited,6=removing,7=stopping).
# TYPE podman_container_state gauge
podman_container_state{id="19286a13dc23",pod_id="",pod_name=""} 2
podman_container_state{id="482113b805f7",pod_id="",pod_name=""} 4
podman_container_state{id="642490688d9c",pod_id="",pod_name=""} 2
podman_container_state{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 5
podman_container_state{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 2

# HELP podman_container_block_input_total Container block input.
# TYPE podman_container_block_input_total counter
podman_container_block_input_total{id="19286a13dc23",pod_id="",pod_name=""} 49152
podman_container_block_input_total{id="482113b805f7",pod_id="",pod_name=""} 0
podman_container_block_input_total{id="642490688d9c",pod_id="",pod_name=""} 1.41533184e+08
podman_container_block_input_total{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_block_input_total{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 0

# HELP podman_container_block_output_total Container block output.
# TYPE podman_container_block_output_total counter
podman_container_block_output_total{id="19286a13dc23",pod_id="",pod_name=""} 1.790976e+06
podman_container_block_output_total{id="482113b805f7",pod_id="",pod_name=""} 8192
podman_container_block_output_total{id="642490688d9c",pod_id="",pod_name=""} 4.69248e+07
podman_container_block_output_total{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_block_output_total{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 0

# HELP podman_container_cpu_seconds_total total CPU time spent for container in seconds.
# TYPE podman_container_cpu_seconds_total counter
podman_container_cpu_seconds_total{id="19286a13dc23",pod_id="",pod_name=""} 83.231904
podman_container_cpu_seconds_total{id="482113b805f7",pod_id="",pod_name=""} 0.069712
podman_container_cpu_seconds_total{id="642490688d9c",pod_id="",pod_name=""} 3.028685
podman_container_cpu_seconds_total{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_cpu_seconds_total{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 0.011687

# HELP podman_container_cpu_system_seconds_total total system CPU time spent for container in seconds.
# TYPE podman_container_cpu_system_seconds_total counter
podman_container_cpu_system_seconds_total{id="19286a13dc23",pod_id="",pod_name=""} 0.007993418
podman_container_cpu_system_seconds_total{id="482113b805f7",pod_id="",pod_name=""} 4.8591e-05
podman_container_cpu_system_seconds_total{id="642490688d9c",pod_id="",pod_name=""} 0.00118734
podman_container_cpu_system_seconds_total{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_cpu_system_seconds_total{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 9.731e-06

# HELP podman_container_created_seconds Container creation time in unixtime.
# TYPE podman_container_created_seconds gauge
podman_container_created_seconds{id="19286a13dc23",pod_id="",pod_name=""} 1.655859887e+09
podman_container_created_seconds{id="482113b805f7",pod_id="",pod_name=""} 1.655859728e+09
podman_container_created_seconds{id="642490688d9c",pod_id="",pod_name=""} 1.655859511e+09
podman_container_created_seconds{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 1.655859858e+09
podman_container_created_seconds{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 1.655859839e+09

# HELP podman_container_started_seconds Container started time in unixtime.
# TYPE podman_container_started_seconds gauge
podman_container_started_seconds{id="19286a13dc23",pod_id="",pod_name=""} 1.659253804e+09
podman_container_started_seconds{id="482113b805f7",pod_id="",pod_name=""} 1.659253804e+09
podman_container_started_seconds{id="642490688d9c",pod_id="",pod_name=""} 1.660642996e+09
podman_container_started_seconds{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 1.66064284e+09
podman_container_started_seconds{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 1.66064284e+09

# HELP podman_container_exit_code Container exit code, if the container has not exited or restarted then the exit code will be 0.
# TYPE podman_container_exit_code gauge
podman_container_exit_code{id="19286a13dc23",pod_id="",pod_name=""} 0
podman_container_exit_code{id="482113b805f7",pod_id="",pod_name=""} 0
podman_container_exit_code{id="642490688d9c",pod_id="",pod_name=""} 0
podman_container_exit_code{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 130
podman_container_exit_code{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 0

# HELP podman_container_exited_seconds Container exited time in unixtime.
# TYPE podman_container_exited_seconds gauge
podman_container_exited_seconds{id="19286a13dc23",pod_id="",pod_name=""} 1.659253805e+09
podman_container_exited_seconds{id="482113b805f7",pod_id="",pod_name=""} 1.659253805e+09
podman_container_exited_seconds{id="642490688d9c",pod_id="",pod_name=""} 1.659253804e+09
podman_container_exited_seconds{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 1.660643511e+09
podman_container_exited_seconds{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 1.660643511e+09

# HELP podman_container_mem_limit_bytes Container memory limit.
# TYPE podman_container_mem_limit_bytes gauge
podman_container_mem_limit_bytes{id="19286a13dc23",pod_id="",pod_name=""} 9.713655808e+09
podman_container_mem_limit_bytes{id="482113b805f7",pod_id="",pod_name=""} 9.713655808e+09
podman_container_mem_limit_bytes{id="642490688d9c",pod_id="",pod_name=""} 9.713655808e+09
podman_container_mem_limit_bytes{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_mem_limit_bytes{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 9.713655808e+09

# HELP podman_container_mem_usage_bytes Container memory usage.
# TYPE podman_container_mem_usage_bytes gauge
podman_container_mem_usage_bytes{id="19286a13dc23",pod_id="",pod_name=""} 1.029062656e+09
podman_container_mem_usage_bytes{id="482113b805f7",pod_id="",pod_name=""} 2.748416e+06
podman_container_mem_usage_bytes{id="642490688d9c",pod_id="",pod_name=""} 3.67616e+07
podman_container_mem_usage_bytes{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_mem_usage_bytes{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 49152

# HELP podman_container_net_input_total Container network input.
# TYPE podman_container_net_input_total counter
podman_container_net_input_total{id="19286a13dc23",pod_id="",pod_name=""} 430
podman_container_net_input_total{id="482113b805f7",pod_id="",pod_name=""} 430
podman_container_net_input_total{id="642490688d9c",pod_id="",pod_name=""} 4323
podman_container_net_input_total{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_net_input_total{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 430

# HELP podman_container_net_output_total Container network output.
# TYPE podman_container_net_output_total counter
podman_container_net_output_total{id="19286a13dc23",pod_id="",pod_name=""} 110
podman_container_net_output_total{id="482113b805f7",pod_id="",pod_name=""} 110
podman_container_net_output_total{id="642490688d9c",pod_id="",pod_name=""} 12071
podman_container_net_output_total{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_net_output_total{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 110

# HELP podman_container_pids Container pid number.
# TYPE podman_container_pids gauge
podman_container_pids{id="19286a13dc23",pod_id="",pod_name=""} 94
podman_container_pids{id="482113b805f7",pod_id="",pod_name=""} 82
podman_container_pids{id="642490688d9c",pod_id="",pod_name=""} 14
podman_container_pids{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 0
podman_container_pids{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 1

# HELP podman_container_rootfs_size_bytes Container root filesystem size in bytes.
# TYPE podman_container_rootfs_size_bytes gauge
podman_container_rootfs_size_bytes{id="19286a13dc23",pod_id="",pod_name=""} 1.452382e+06
podman_container_rootfs_size_bytes{id="482113b805f7",pod_id="",pod_name=""} 1.135744e+06
podman_container_rootfs_size_bytes{id="642490688d9c",pod_id="",pod_name=""} 1.72771905e+08
podman_container_rootfs_size_bytes{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 1.135744e+06
podman_container_rootfs_size_bytes{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 1.035744e+06

# HELP podman_container_rw_size_bytes Container top read-write layer size in bytes.
# TYPE podman_container_rw_size_bytes gauge
podman_container_rw_size_bytes{id="19286a13dc23",pod_id="",pod_name=""} 0
podman_container_rw_size_bytes{id="482113b805f7",pod_id="",pod_name=""} 0
podman_container_rw_size_bytes{id="642490688d9c",pod_id="",pod_name=""} 26261
podman_container_rw_size_bytes{id="ad36e85960a1",pod_id="3e8bae64e9af",pod_name="pod01"} 3551
podman_container_rw_size_bytes{id="dda983cc3ecf",pod_id="3e8bae64e9af",pod_name="pod01"} 0
```

#### `pod`

```shell
# HELP podman_pod_state Pods current state current state (-1=unknown,0=created,1=error,2=exited,3=paused,4=running,5=degraded,6=stopped).
# TYPE podman_pod_state gauge
podman_pod_state{id="3e8bae64e9af"} 5
podman_pod_state{id="959a0a3530db"} 0
podman_pod_state{id="d05cda23085a"} 2

# HELP podman_pod_info Pod information
# TYPE podman_pod_info gauge
podman_pod_info{id="3e8bae64e9af",infra_id="dda983cc3ecf",name="pod01"} 1
podman_pod_info{id="959a0a3530db",infra_id="22e3d69be889",name="pod02"} 1
podman_pod_info{id="d05cda23085a",infra_id="390ac740fa80",name="pod03"} 1

# HELP podman_pod_containers Number of containers in a pod.
# TYPE podman_pod_containers gauge
podman_pod_containers{id="3e8bae64e9af"} 2
podman_pod_containers{id="959a0a3530db"} 1
podman_pod_containers{id="d05cda23085a"} 1

# HELP podman_pod_created_seconds Pods creation time in unixtime.
# TYPE podman_pod_created_seconds gauge
podman_pod_created_seconds{id="3e8bae64e9af"} 1.655859839e+09
podman_pod_created_seconds{id="959a0a3530db"} 1.655484892e+09
podman_pod_created_seconds{id="d05cda23085a"} 1.655489348e+09
```

#### `image`

```shell
# HELP podman_image_info Image information.
# TYPE podman_image_info gauge
podman_image_info{id="48565a8e6250",parent_id="",repository="docker.io/bitnami/prometheus",tag="latest",digest="sha256:4d7fdebe2a853aceb15019554b56e58055f7a746c0b4095eec869d5b6c11987e"} 1
podman_image_info{id="62aedd01bd85",parent_id="",repository="docker.io/library/busybox",tag="latest",digest="sha256:6d9ac9237a84afe1516540f40a0fafdc86859b2141954b4d643af7066d598b74"} 1
podman_image_info{id="75c013514322",parent_id="",repository="docker.io/library/sonarqube",tag="latest",digest="sha256:548f3d4246cda60c311a035620c26ea8fb21b3abc870c5806626a32ef936982b"} 1
podman_image_info{id="a45fa0117c2b",parent_id="",repository="localhost/podman-pause",tag="4.1.0-1651853754",digest="sha256:218169c5590870bb95c06e9f7e80ded58f6644c1974b0ca7f2c3405b74fc3b57"} 1
podman_image_info{id="b260a49eebf9",parent_id="",repository="docker.io/library/httpd",tag="latest",digest="sha256:ba846154ade27292d216cce2d21f1c7e589f3b66a4a643bff0cdd348efd17aa3"} 1
podman_image_info{id="c4b778290339",parent_id="b260a49eebf9",repository="docker.io/grafana/grafana",tag="latest",digest="sha256:7567a7c70a3c1d75aeeedc968d1304174a16651e55a60d1fb132a05e1e63a054"} 1

# HELP podman_image_created_seconds Image creation time in unixtime.
# TYPE podman_image_created_seconds gauge
podman_image_created_seconds{id="48565a8e6250",repository="docker.io/bitnami/prometheus",tag="latest"} 1.655436988e+09
podman_image_created_seconds{id="62aedd01bd85",repository="docker.io/library/busybox",tag="latest"} 1.654651161e+09
podman_image_created_seconds{id="75c013514322",repository="docker.io/library/sonarqube",tag="latest"} 1.654883091e+09
podman_image_created_seconds{id="a45fa0117c2b",repository="localhost/podman-pause",tag="4.1.0-1651853754"} 1.655484887e+09
podman_image_created_seconds{id="b260a49eebf9",repository="docker.io/library/httpd",tag="latest"} 1.655163309e+09
podman_image_created_seconds{id="c4b778290339",repository="docker.io/grafana/grafana",tag="latest"} 1.655132996e+09

# HELP podman_image_size Image size
# TYPE podman_image_size gauge
podman_image_size{id="48565a8e6250",repository="docker.io/bitnami/prometheus",tag="latest"} 5.11822059e+08
podman_image_size{id="62aedd01bd85",repository="docker.io/library/busybox",tag="latest"} 1.468102e+06
podman_image_size{id="75c013514322",repository="docker.io/library/sonarqube",tag="latest"} 5.35070053e+08
podman_image_size{id="a45fa0117c2b",repository="localhost/podman-pause",tag="4.1.0-1651853754"} 815742
podman_image_size{id="b260a49eebf9",repository="docker.io/library/httpd",tag="latest"} 1.49464899e+08
podman_image_size{id="c4b778290339",repository="docker.io/grafana/grafana",tag="latest"} 2.98969093e+08
```

#### `network`

```shell
# HELP podman_network_info Network information.
# TYPE podman_network_info gauge
podman_network_info{driver="bridge",id="2f259bab93aa",interface="podman0",labels="",name="podman"} 1
podman_network_info{driver="bridge",id="420272a98a4c",interface="podman3",labels="",name="network03"} 1
podman_network_info{driver="bridge",id="6eb310d4b0bb",interface="podman2",labels="",name="network02"} 1
podman_network_info{driver="bridge",id="a5a6391121a5",interface="podman1",labels="",name="network01"} 1
```

#### `volume`

```shell
# HELP podman_volume_info Volume information.
# TYPE podman_volume_info gauge
podman_volume_info{driver="local",mount_point="/home/navid/.local/share/containers/storage/volumes/vol01/_data",name="vol01"} 1
podman_volume_info{driver="local",mount_point="/home/navid/.local/share/containers/storage/volumes/vol02/_data",name="vol02"} 1
podman_volume_info{driver="local",mount_point="/home/navid/.local/share/containers/storage/volumes/vol03/_data",name="vol03"} 1

# HELP podman_volume_created_seconds Volume creation time in unixtime.
# TYPE podman_volume_created_seconds gauge
podman_volume_created_seconds{name="vol01"} 1.655484915e+09
podman_volume_created_seconds{name="vol02"} 1.655484926e+09
podman_volume_created_seconds{name="vol03"} 1.65548493e+09
```

#### `system`

```shell
# HELP podman_system_api_version Podman system api version.
# TYPE podman_system_api_version gauge
podman_system_api_version{version="4.1.1"} 1

# HELP podman_system_buildah_version Podman system buildahVer version.
# TYPE podman_system_buildah_version gauge
podman_system_buildah_version{version="1.26.1"} 1

# HELP podman_system_conmon_version Podman system conmon version.
# TYPE podman_system_conmon_version gauge
podman_system_conmon_version{version="2.1.0"} 1

# HELP podman_system_runtime_version Podman system runtime version.
# TYPE podman_system_runtime_version gauge
podman_system_runtime_version{version="crun version 1.4.5"} 1
```

## License

Licensed under the [Apache 2.0](LICENSE) license.
