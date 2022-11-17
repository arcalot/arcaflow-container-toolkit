FROM golang:1.18 AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./carpenter.go


FROM quay.io/centos/centos:stream8
RUN dnf -y install dnf-plugins-core &&\
    dnf -y install golang &&\
    dnf -y config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo &&\
    dnf -y install docker-ce-cli
RUN dnf -y module install python39 &&\
    dnf -y install python39 python39-pip &&\
    python3.9 -m pip install --user --upgrade flake8

COPY --from=builder /build/carpenter /
COPY .carpenter.yaml /.carpenter.yaml
WORKDIR /


ENTRYPOINT ["/carpenter"]
CMD ["build"]


LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-image-builder"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Plugin Image Builder"