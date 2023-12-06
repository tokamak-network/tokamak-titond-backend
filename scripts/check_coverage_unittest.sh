#!/bin/bash

expected=70
actual=$(go tool cover -func=coverage.out | grep ^total | awk '{print $3}' | cut -d. -f1)
echo "Actual coverage is $actual%"
if [ "$actual" -lt "$expected" ]; then
  echo "Error: The coverage is $actual%. The threshold is $expected%."
  exit 1
fi
