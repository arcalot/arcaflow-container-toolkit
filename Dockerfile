FROM quay.io/centos/centos:stream8
RUN dnf -y install dnf-plugins-core &&\
    dnf -y config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo &&\
    dnf -y install golang docker-ce-cli git make
RUN dnf -y module install python39 &&\
    dnf -y install python39 python39-pip &&\
    python3.9 -m pip install --user --upgrade flake8


RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.50.1 &&\
    git clone https://github.com/GoTestTools/limgo.git
WORKDIR /limgo
RUN make build

COPY . /build
WORKDIR /build

RUN go test ./... -coverprofile=cov.out &&\
    ../limgo/limgo -coverfile=cov.out -config=.limgo.json -outfile=/limgo_cov.txt

RUN CGO_ENABLED=0 go build ./carpenter.go
RUN mv carpenter /carpenter
COPY .carpenter.yaml /.carpenter.yaml
WORKDIR /

ENTRYPOINT ["/carpenter"]
CMD ["build"]

LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-plugin-image-builder"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Plugin Image Builder"