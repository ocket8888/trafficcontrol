#!/bin/sh -l
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.

set -e

GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"
cd "$SRCDIR/trafficcontrol"

FILES="$(/usr/local/go/bin/go fmt $INPUT_DIR)"
if [ -z "$FILES" ]; then
	exit 0
fi
echo "$FILES"

for f in "$FILES"; do
	echo "$f" >&2
done

git --no-pager diff >&2
exit 1
