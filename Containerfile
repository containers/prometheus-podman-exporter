FROM --platform=${BUILDPLATFORM:-linux/amd64} quay.io/prometheus/busybox:latest
LABEL maintainer="Navid Yaghoobi <navidys@fedoraproject.org>"

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

COPY ./bin/remote/prometheus-podman-exporter-${TARGETARCH} /bin/podman_exporter

EXPOSE 9882
USER nobody
ENTRYPOINT [ "/bin/podman_exporter" ]
