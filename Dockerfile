# Build stage
FROM golang:1.22@sha256:cefea7fa6852b85f0042ce9d4b883c7e0b03b2bcb25972372d59e4f7c4367c04 AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./act.go

# Main stage
FROM python:3.12.2-slim-bullseye@sha256:6229a3e2fd625bb99efb87e38e7ca4aff6351636ce41b3f229f672b18be8134c
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
