# Build stage
FROM golang:1.21@sha256:b113af1e8b06f06a18ad41a6b331646dff587d7a4cf740f4852d16c49ed8ad73 AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./act.go

# Main stage
FROM python:3.11.5-slim-bullseye@sha256:9f35f3a6420693c209c11bba63dcf103d88e47ebe0b205336b5168c122967edf
RUN python -m ensurepip
RUN python -m pip install --user --upgrade flake8

COPY --from=builder /build/act /
COPY .act.yaml /.act.yaml
WORKDIR /

ENTRYPOINT ["/act"]
CMD ["build"]


LABEL org.opencontainers.image.source="https://github.com/arcalot/arcaflow-container-toolkit"
LABEL org.opencontainers.image.licenses="Apache-2.0+GPL-2.0-only"
LABEL org.opencontainers.image.vendor="Arcalot project"
LABEL org.opencontainers.image.authors="Arcalot contributors"
LABEL org.opencontainers.image.title="Arcaflow Container Toolkit"
