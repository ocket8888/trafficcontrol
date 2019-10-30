# go-test Docker action

This action runs Go unit tests.

## Inputs

### `dir`

**Required** Directory in which to run tests.

## Outputs

### `result`

Test results.

### `exit-code`

Exit code of the test command

## Example usage

uses: actions/go-test@v1
with:
  dir: './lib/...'
