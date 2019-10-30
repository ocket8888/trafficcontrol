#!/bin/sh -l

output=$(go test $@ 2>&1)
code=$?

echo ::set-output name=result::$output
echo ::set-output name=exit-code::$code
