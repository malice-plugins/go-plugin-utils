#!/bin/sh

set -e

GOLANG_VERSION=1.7.5
GOLANG_SRC_URL=https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
GOLANG_SRC_SHA256=4e834513a2079f8cbbd357502cccaac9507fd00a1efe672375798858ff291815
PATCH_URL=https://raw.githubusercontent.com/maliceio/go-plugin-utils/master/scripts/no-pic.patch

echo "Upgrade to Golang $GOLANG_VERSION..."

export GOROOT_BOOTSTRAP="$(go env GOROOT)"

wget -q "$PATCH_URL" -O /no-pic.patch
wget -q "$GOLANG_SRC_URL" -O golang.tar.gz
echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c -

tar -C /usr/local -xzf golang.tar.gz
rm golang.tar.gz

cd /usr/local/go/src

patch -p2 -i /no-pic.patch
./make.bash
rm -rf /*.patch
