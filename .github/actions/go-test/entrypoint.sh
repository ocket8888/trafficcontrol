#!/bin/sh -l

set -e

if [ -z "$INPUT_DIR" ]; then
	# There's a bug in "defaults" for inputs
	INPUT_DIR="./lib/..."
fi

GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"
cd "$SRCDIR/trafficcontrol"

# Need to fetch golang.org/x/* dependencies
/usr/local/go/bin/go get -v $INPUT_DIR
/usr/local/go/bin/go test -v $INPUT_DIR
