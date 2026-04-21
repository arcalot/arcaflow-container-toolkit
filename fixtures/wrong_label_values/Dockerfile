FROM quay.io/centos/centos:stream8

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

LABEL org.opencontainers.image.source="https://github.com/example/example-plugin"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.vendor="Example project"
LABEL org.opencontainers.image.authors="Example contributors"
LABEL org.opencontainers.image.title="Example Plugin"
LABEL io.github.arcalot.arcaflow.plugin.version="1"
