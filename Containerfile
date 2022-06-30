ARG ARCH="amd64"
ARG OS="linux"
FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest

LABEL maintainer="Navid Yaghoobi <navidys@fedoraproject.org>"

COPY ./bin/remote/prometheus-podman-exporter /bin/podman_exporter

EXPOSE 9882
USER nobody
ENTRYPOINT [ "/bin/podman_exporter" ]
