# Build stage
FROM golang:1.22@sha256:34ce21a9696a017249614876638ea37ceca13cdd88f582caad06f87a8aa45bf3 AS builder
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./act.go

# Main stage
FROM python:3.12.2-slim-bullseye@sha256:1c0da9b35e7fbba5441c9652e93194450f849ba69599dee38ebf9a04c011dc42
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
