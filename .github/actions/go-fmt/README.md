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

# go-fmt Docker action
This action runs `gofmt` on all Go source files (including test files) under the provided directory. It exits with a failure if one or more files changed after formatting.

## Inputs

### `dir`
**Required** Directory in which to look for Go source files

### `exit-code`
1 if there were files not properly formatted, 0 if everything ran fine and no changes needed to be made.

## Example usage
```yaml
uses: actions/go-fmt
with:
  dir: './lib/...'
```