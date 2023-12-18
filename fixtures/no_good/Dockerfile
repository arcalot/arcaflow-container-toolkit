FROM quay.io/centos/centos:stream8@sha256:34aaf8788a2467f602c5772884448236bb41dfe1691a78dee33053bb24474395
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