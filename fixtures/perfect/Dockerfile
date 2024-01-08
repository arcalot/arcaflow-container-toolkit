FROM quay.io/centos/centos:stream8@sha256:a8692b39e546eed9177d495db1edfd97bb6de70b9527f58aeb72f90b687c3426

RUN dnf -y module install python39 && dnf -y install python39 python39-pip
RUN dnf -y install fio-3.19-3.el8
RUN mkdir /plugin
ADD https://raw.githubusercontent.com/arcalot/arcaflow-plugins/main/LICENSE /plugin
COPY requirements.txt /plugin
RUN pip3 install --requirement /plugin/requirements.txt
COPY fio_plugin.py /plugin
COPY fio_schema.py /plugin
COPY test_fio_plugin.py /plugin
COPY fixtures /plugin/fixtures

WORKDIR /plugin
RUN python3.9 test_fio_plugin.py

ENTRYPOINT ["python3.9", "fio_plugin.py" ]
CMD []

LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-fio"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Fio Arcalot Plugin"
LABEL io.github.arcalot.arcaflow.plugin.version="1"