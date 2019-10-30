#!/bin/sh -l

set -e

echo "Called with $# args: $@"

DIR="$1"

if [ -z "$DIR" ]; then
	# stupid inputs don't stupid work
	DIR="./lib/..."
fi

GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"
cd "$SRCDIR/trafficcontrol"

/usr/local/go/bin/go test -v $DIR
