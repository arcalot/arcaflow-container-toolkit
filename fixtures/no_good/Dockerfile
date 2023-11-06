FROM quay.io/centos/centos:stream8@sha256:e7f228fe74eeac927a3133ae78a75aac1f28f6dff284616a7b2b10b5769b6275
RUN dnf -y install dnf-plugins-core &&\
    dnf -y config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo &&\
    dnf -y install golang docker-ce-cli
RUN dnf -y module install python39 &&\
    dnf -y install python39 python39-pip &&\
    python3.9 -m pip install --user --upgrade flake8

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./act.go
RUN mv act /act
COPY .act.yaml /.act.yaml
WORKDIR /

ENTRYPOINT ["/act"]
CMD ["build"]

LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-container-toolkit"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Plugin Image Builder"