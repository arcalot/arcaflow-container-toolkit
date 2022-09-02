FROM quay.io/centos/centos:stream8

RUN dnf -y module install python39 && dnf -y install python39 python39-pip
RUN mkdir /plugin
ADD https://raw.githubusercontent.com/arcalot/arcaflow-plugins/main/LICENSE /
RUN python3.9 -m pip install --user --upgrade pycodestyle flake8


WORKDIR /plugin
ENTRYPOINT ["python3", "-m", "flake8", "--show-source", "." ]

LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-image-builder/containerfiles/python"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Python Plugin Image Builder"
