## Installation Guide

- [**Building From Source**](#building-from-source)
- [**Container Image**](#container-image)
- [**Installing Packaged Versions**](#installing-packaged-versions)
  - [**Centos Stream**](#centos-stream)
  - [**Enterprise Linux (EPEL)**](#enterprise-linux-epel)
  - [**Fedora**](#fedora)

## Building From Source

prometheus-podman-exporter is using go v1.17 or above.

1. Clone the repo
2. Install dependencies
    * Fedora

        ```shell
        $ sudo dnf install -y btrfs-progs-devel device-mapper-devel gpgme-devel libassuan-devel
        ```

    * Debian

        ```shell
        $ sudo apt-get -y install libgpgme-dev libbtrfs-dev libdevmapper-dev libassuan-dev pkg-config
        ```

2. Build and run the executable

    ```shell
    $ make binary
    $ ./bin/prometheus-podman-exporter
    ```
## Container Image

```shell
podman run -e CONTAINER_HOST=tcp://<host>:<port> --network=host -p 9882:9882 quay.io/navidys/prometheus-podman-exporter:latest
```

## Installing Packaged Versions

### Centos Stream

RPM package is available through [COPR repo](https://copr.fedorainfracloud.org/coprs/navidys/prometheus-podman-exporter/).

### Enterprise Linux (EPEL)

RPM package is available through [COPR repo](https://copr.fedorainfracloud.org/coprs/navidys/prometheus-podman-exporter/).

### Fedora

RPM package is available through [COPR repo](https://copr.fedorainfracloud.org/coprs/navidys/prometheus-podman-exporter/).

```
$ sudo dnf copr enable navidys/prometheus-podman-exporter
$ sudo dnf install prometheus-podman-exporter
```
