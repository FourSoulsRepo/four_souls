#!/usr/bin/env bash
# Builds the dedicated server archives of a release (8.17, N-19).
#
#   release/server.sh <version> [out-dir]
#
# The server is plain Go: CGO_ENABLED=0 makes the Linux binaries static,
# so they run on any distribution, Alpine included. Names carry no
# version, so releases/latest/download/<name> always finds the newest
# one; the Dockerfile relies on it. macOS gets one universal binary,
# which needs lipo, so this runs on macOS.
set -euo pipefail

version=${1:?usage: release/server.sh <version> [out-dir]}
out=${2:-dist}
root=$(cd "$(dirname "$0")/.." && pwd)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
ldflags="-X github.com/FourSoulsRepo/four_souls/internal/version.App=$version"
export CGO_ENABLED=0
# No ._ metadata files in the archives.
export COPYFILE_DISABLE=1

mkdir -p "$out"
out=$(cd "$out" && pwd)
cd "$root"

# build <goos> <goarch> <file>
build() {
	# The server needs card text, not card images.
	GOOS=$1 GOARCH=$2 go build -tags embed -ldflags "$ldflags" -o "$3" ./cmd/server
}

# pack <name> <tar|zip>
pack() {
	cp release/SERVER.txt "$work/$1/README.txt"
	cp LICENSE NOTICE "$work/$1/"
	case $2 in
	tar) tar -czf "$out/$1.tar.gz" -C "$work" "$1" ;;
	zip) (cd "$work" && zip -qrX "$out/$1.zip" "$1") ;;
	esac
}

for arch in amd64 arm64; do
	name=four-souls-server-linux-$arch
	mkdir "$work/$name"
	build linux "$arch" "$work/$name/four_souls_server"
	pack "$name" tar

	name=four-souls-server-windows-$arch
	mkdir "$work/$name"
	build windows "$arch" "$work/$name/four_souls_server.exe"
	pack "$name" zip
done

name=four-souls-server-macos-universal
mkdir "$work/$name"
build darwin amd64 "$work/server-amd64"
build darwin arm64 "$work/server-arm64"
lipo -create -output "$work/$name/four_souls_server" "$work/server-amd64" "$work/server-arm64"
pack "$name" zip

ls -la "$out"
