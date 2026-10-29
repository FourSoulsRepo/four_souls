# Four Souls dedicated server from the latest GitHub release (8.19).
#
# This file builds no source code. Put it and build.sh in an empty
# folder, then:
#
#   ./build.sh                     # image four-souls-server:0.3.2
#   docker run -d --name four-souls --restart unless-stopped \
#     -p 4774:4774 -v four-souls-data:/data four-souls-server:0.3.2
#
# Update: run ./build.sh again. It passes the newest version, so Docker
# downloads only when a new release is out; then replace the container
# (wiki page "Server in Docker"). Without VERSION, a cached download is
# kept: add --no-cache to fetch the latest again.
#
# A fixed release instead of the latest: ./build.sh 0.3.1
# The image shows its version in its label and in the server's first
# log line ("app 0.3.2").
# Settings: FOUR_SOULS_* ENV variables, e.g. -e FOUR_SOULS_NAME=Basement
# (wiki page "Dedicated server on a VPS", Options).

ARG ALPINE=3.24.2

FROM alpine:${ALPINE} AS download
ARG VERSION=latest
# Set by BuildKit: amd64 or arm64. The legacy builder leaves it empty,
# so the CPU is read from uname then.
ARG TARGETARCH
RUN set -eu; \
    arch=${TARGETARCH:-$(uname -m)}; \
    case "$arch" in \
      amd64|x86_64) arch=amd64 ;; \
      arm64|aarch64) arch=arm64 ;; \
      *) echo "unsupported CPU: $arch" >&2; exit 1 ;; \
    esac; \
    name="four-souls-server-linux-$arch"; \
    url="https://github.com/FourSoulsRepo/four_souls/releases/latest/download/$name.tar.gz"; \
    if [ "$VERSION" != latest ]; then \
      url="https://github.com/FourSoulsRepo/four_souls/releases/download/v${VERSION#v}/$name.tar.gz"; \
    fi; \
    wget -q -O /tmp/server.tar.gz "$url"; \
    tar -xzf /tmp/server.tar.gz -C /tmp; \
    install -m 755 "/tmp/$name/four_souls_server" /four_souls_server

FROM alpine:${ALPINE}
ARG VERSION=latest
LABEL org.opencontainers.image.title="Four Souls server" \
    org.opencontainers.image.version="${VERSION}"
RUN addgroup -S four-souls \
    && adduser -S -D -H -G four-souls -h /data four-souls \
    && mkdir -p /data \
    && chown four-souls:four-souls /data
COPY --from=download /four_souls_server /usr/local/bin/four_souls_server
ENV FOUR_SOULS_PORT=4774 \
    FOUR_SOULS_RECORDS=/data/records \
    FOUR_SOULS_SAVES=/data/saves
WORKDIR /data
VOLUME /data
EXPOSE 4774
USER four-souls
# /healthz answers 200 while the server's lobby answers.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O /dev/null "http://127.0.0.1:${FOUR_SOULS_PORT}/healthz" || exit 1
# SIGTERM (docker stop) closes every game cleanly, within 10 seconds.
STOPSIGNAL SIGTERM
ENTRYPOINT ["four_souls_server"]
