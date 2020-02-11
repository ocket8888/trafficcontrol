<!--
	Licensed to the Apache Software Foundation (ASF) under one
	or more contributor license agreements.  See the NOTICE file
	distributed with this work for additional information
	regarding copyright ownership.  The ASF licenses this file
	to you under the Apache License, Version 2.0 (the
	"License"); you may not use this file except in compliance
	with the License.  You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing,
	software distributed under the License is distributed on an
	"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
	KIND, either express or implied.  See the License for the
	specific language governing permissions and limitations
	under the License.
-->
# convert_profile
`convert_profile` is a tool for converting ATC Profiles meant for use with ATS cache servers using older
versions of ATS to use newer options compatible with newer versions of ATS. Specifically, a ruleset is
given in [convert622to713.json](convert622to713.json) (and in YAML format in
[convert622to713.yaml](convert622to713.yaml)) for converting Profiles made for version 6.22 (or possibly
lower) to use newer options compatible with ATS 7.13.

## Building
`convert_profile` is built using Go, so to build just use
`go build github.com/apache/trafficcontrol/tools/convert_profile`.

## Usage
```bash
convert_profile -i IN -r RULES [-o OUT] [-f]
```

### Options

#### -i/--input_profile IN
Reads the file `IN` as an ATC Profile object encoded in JSON. This option is not optional; there must be
input to process.

#### -r/--rules RULES
Reads the file `RULES` as a set of rules to use when converting Profiles. It may be encoded in JSON or
YAML. This option is not optional; there must be rules to use when converting. (For more information
refer to the
[official documentation](https://traffic-control-cdn.readthedocs.io/en/latest/tools/convert_profile.html))

#### -o/--out OUT
Writes the resulting Profile object to the file given by `OUT`. If not given, the result is printed to
stdout.

#### -f/--force
If given, existing values in the input Profile will be "clobbered", replacing them with the
recommendations given by the rules even if the value doesn't match a pattern given by the rule set.

## Testing
Unit testing is available via `go test`.
