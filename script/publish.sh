#!/bin/sh
set -e

latest_tag() {
  git tag -l | egrep "^v[0-9]+\.[0-9]+\.[0-9]+" | cut -d"-" -f 1 | sort | tail -n 1
}

publish_release() {
  echo "Triggering pkg.go.dev about new smtpmock release..."
  curl -X POST "https://pkg.go.dev/fetch/github.com/mocktools/go-smtp-mock@$(latest_tag)"
}

publish_release
