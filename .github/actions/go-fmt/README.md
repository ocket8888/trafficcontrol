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
