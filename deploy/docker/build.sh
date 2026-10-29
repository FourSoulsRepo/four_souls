#!/bin/sh
# Builds the server image with two tags: its app version, e.g. four-souls-server:0.3.2,
# and latest, which docker-compose.yml uses. Without an argument it takes the newest
# release; give a version to pin another one.
#
#   ./build.sh          # newest release
#   ./build.sh 0.3.1    # that release
set -eu

version=${1:-}
if [ -z "$version" ]; then
	version=$(curl -fsSI -o /dev/null -w '%{redirect_url}' \
		https://github.com/FourSoulsRepo/four_souls/releases/latest | sed 's|.*/tag/v||')
fi
version=${version#v}

cd "$(dirname "$0")"
docker build \
	--build-arg VERSION="$version" \
	-t "four-souls-server:$version" \
	-t four-souls-server:latest \
	.
echo "built four-souls-server:$version and four-souls-server:latest"
