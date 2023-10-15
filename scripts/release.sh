#!/usr/bin/env bash

TAG=$(grep "const VERSION string" < version/main.go | awk -F'"' '{$0=$2}1')

git push --delete origin "v${TAG}"
git tag --delete "v${TAG}"
git tag "v${TAG}"
git push --tags --force
