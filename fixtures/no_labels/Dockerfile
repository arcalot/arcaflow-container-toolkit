FROM quay.io/centos/centos:stream8

RUN dnf -y module install python39 && dnf -y install python39 python39-pip
RUN mkdir /plugin
ADD https://raw.githubusercontent.com/arcalot/arcaflow-plugins/main/LICENSE /plugin
COPY requirements.txt /plugin
COPY fio_plugin.py /plugin

WORKDIR /plugin

ENTRYPOINT ["python3.9", "fio_plugin.py" ]
CMD []
