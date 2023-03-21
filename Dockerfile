FROM golang:1.18 AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./act.go


FROM quay.io/centos/centos:stream8
RUN dnf -y module install python39 &&\
    dnf -y install python39 python39-pip &&\
    python3.9 -m pip install --user --upgrade flake8

COPY --from=builder /build/act /
COPY .act.yaml /.act.yaml
WORKDIR /


ENTRYPOINT ["/act"]
CMD ["build"]


LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-image-builder"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Plugin Image Builder"