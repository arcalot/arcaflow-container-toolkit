FROM quay.io/centos/centos:stream8

RUN mkdir /plugin
ADD https://raw.githubusercontent.com/arcalot/arcaflow-plugins/main/LICENSE /

RUN dnf -y module install python39 && dnf -y install python39 python39-pip &&\
    python3.9 -m pip install --user --upgrade flake8
RUN curl -sSL https://git.io/g-install | sh -s &&\
    g install 1.17


WORKDIR /plugin
# ENTRYPOINT ["python3", "-m", "flake8", "--show-source", "." ]

LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-image-builder"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Plugin Image Builder"