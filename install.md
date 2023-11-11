## Installation Guide

- [**Building From Source**](#building-from-source)
- [**Container Image**](#container-image)
- [**Installing Packaged Versions**](#installing-packaged-versions)
  - [**Centos Stream**](#centos-stream)
  - [**Fedora**](#fedora)
  - [**RHEL**](#rhel)
  

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

* Using unix socket (rootless):

    ```shell
    $ systemctl start --user podman.socket
    $ podman run -e CONTAINER_HOST=unix:///run/podman/podman.sock -v $XDG_RUNTIME_DIR/podman/podman.sock:/run/podman/podman.sock --userns=keep-id --security-opt label=disable quay.io/navidys/prometheus-podman-exporter
    ```

* Using unix socket (root):

    ```
    # systemctl start podman.socket
    # podman run -e CONTAINER_HOST=unix:///run/podman/podman.sock -v /run/podman/podman.sock:/run/podman/podman.sock -u root --security-opt label=disable quay.io/navidys/prometheus-podman-exporter
    ```

* Using TCP:

    ```shell
    $ podman system service --time=0 tcp://<ip>:<port>
    $ podman run -e CONTAINER_HOST=tcp://<ip>:<port> --network=host -p 9882:9882 quay.io/navidys/prometheus-podman-exporter:latest
    ```

## Installing Packaged Versions

### Centos Stream

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```
$ sudo dnf install prometheus-podman-exporter
```

### Fedora

```
$ sudo dnf install prometheus-podman-exporter
```

### RHEL

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```
$ sudo dnf install prometheus-podman-exporter
```