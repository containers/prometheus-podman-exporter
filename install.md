## Installation Guide

- [**Building From Source**](#building-from-source)
- [**Container Image**](#container-image)
- [**Installing Packaged Versions**](#installing-packaged-versions)
  - [**AlmaLinux, Rocky Linux**](#almalinux-rocky-linux)
  - [**Arch Linux (AUR)**](#arch-linux-aur)
  - [**Centos Stream**](#centos-stream)
  - [**Fedora**](#fedora)
  - [**RHEL**](#rhel)
  - [**Gentoo**](#gentoo)


## Building From Source

prometheus-podman-exporter is using go v1.17 or above.

1. Clone the repo
2. Install dependencies

    ```shell
    $ sudo dnf -y install btrfs-progs-devel device-mapper-devel gpgme-devel libassuan-devel
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
    $ podman run -e CONTAINER_HOST=unix:///run/podman/podman.sock -v $XDG_RUNTIME_DIR/podman/podman.sock:/run/podman/podman.sock -p 9882:9882 --userns=keep-id --security-opt label=disable quay.io/navidys/prometheus-podman-exporter
    ```

* Using unix socket (root):

    ```
    # systemctl start podman.socket
    # podman run -e CONTAINER_HOST=unix:///run/podman/podman.sock -v /run/podman/podman.sock:/run/podman/podman.sock -u root -p 9882:9882 --security-opt label=disable quay.io/navidys/prometheus-podman-exporter
    ```

* Using TCP:

    ```shell
    $ podman system service --time=0 tcp://<ip>:<port>
    $ podman run -e CONTAINER_HOST=tcp://<ip>:<port> --network=host -p 9882:9882 quay.io/navidys/prometheus-podman-exporter:latest
    ```

## Installing Packaged Versions

### AlmaLinux, Rocky Linux

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```shell
$ sudo dnf -y install prometheus-podman-exporter
```

### Arch Linux (AUR)

```shell
$ yay -S prometheus-podman-exporter
```

### Centos Stream

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```shell
$ sudo dnf -y install prometheus-podman-exporter
```

### Fedora

```shell
$ sudo dnf -y install prometheus-podman-exporter
```

### RHEL

Enable [EPEL repository](https://docs.fedoraproject.org/en-US/epel/) and then run:

```shell
$ sudo dnf -y install prometheus-podman-exporter
```

### Gentoo

```shell
$ sudo emerge app-metrics/prometheus-podman-exporter
```
