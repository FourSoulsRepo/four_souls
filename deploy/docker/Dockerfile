# Four Souls dedicated server from the latest GitHub release (8.19).
#
# This file is all you need; it builds no source code. Put it in an
# empty folder, then:
#
#   docker build -t four-souls-server .
#   docker run -d --name four-souls --restart unless-stopped \
#     -p 4774:4774 -v four-souls-data:/data four-souls-server
#
# Update: run the same docker build again. Docker checks the release
# link on every build and downloads only when a new release is out;
# then replace the container (wiki page "Server in Docker").
#
# A fixed release instead of the latest: --build-arg VERSION=1.2.3
# Settings: FOUR_SOULS_* ENV variables, e.g. -e FOUR_SOULS_NAME=Basement
# (wiki page "Dedicated server on a VPS", Options).

ARG ALPINE=3.24.2

FROM alpine:${ALPINE} AS download
ARG VERSION=latest
# Set by Docker: amd64 or arm64.
ARG TARGETARCH
# Always the latest server: Docker fetches this link on every build, so
# a new release reruns the steps below and an old one keeps the cache.
ADD https://github.com/FourSoulsRepo/four_souls/releases/latest/download/four-souls-server-linux-${TARGETARCH}.tar.gz /tmp/latest.tar.gz
RUN set -eu; \
    name="four-souls-server-linux-${TARGETARCH}"; \
    archive=/tmp/latest.tar.gz; \
    if [ "$VERSION" != latest ]; then \
      archive=/tmp/pinned.tar.gz; \
      wget -q -O "$archive" "https://github.com/FourSoulsRepo/four_souls/releases/download/v${VERSION#v}/$name.tar.gz"; \
    fi; \
    tar -xzf "$archive" -C /tmp; \
    install -m 755 "/tmp/$name/four_souls_server" /four_souls_server

FROM alpine:${ALPINE}
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
