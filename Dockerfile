FROM golang:1.21 as builder
WORKDIR /src
ARG VERSION

RUN apt-get update && apt-get upgrade -y && apt-get install -y make gcc g++ ca-certificates

COPY . .

RUN VERSION=${VERSION} make build

FROM debian:stable-slim AS runtime
LABEL maintainer="Moov <oss@moov.io>"

WORKDIR /

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/bin/bai2 /app/
COPY /configs/ /configs/

ENV HTTP_PORT=8484
ENV HEALTH_PORT=9494

EXPOSE ${HTTP_PORT}/tcp
EXPOSE ${HEALTH_PORT}/tcp

HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
	CMD curl -f http://localhost:${HEALTH_PORT}/live || exit 1

VOLUME [ "/data", "/configs" ]

ENTRYPOINT ["/app/bai2"]
