#!/bin/sh -l

set -e

GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"
cd "$SRCDIR/trafficcontrol"

files="$(/usr/local/go/bin/go fmt $INPUT_DIR)"
if [ -z "$FILES" ]; then
	exit 0
fi

for f in "$files"; do
	echo "$f" >&2
done

git --no-pager diff >&2
exit 1
