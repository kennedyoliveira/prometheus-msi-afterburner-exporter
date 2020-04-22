# This is for local usage, the builds pushed to the registry will be built with Dockerfile-cross and buildx
FROM ubuntu
COPY ./bin/prometheus-msi-afterburner-exporter /bin/prometheus-msi-afterburner-exporter

ENTRYPOINT ["/bin/prometheus-msi-afterburner-exporter"]