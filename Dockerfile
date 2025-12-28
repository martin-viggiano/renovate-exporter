FROM alpine:3

COPY renovate-exporter /

ENTRYPOINT ["/renovate-exporter"]