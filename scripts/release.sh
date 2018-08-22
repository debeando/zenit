#!/usr/bin/env bash

if ! type "jq" > /dev/null; then
  echo "Require command tool, please install: jq"
  exit 1
fi

if [ -z "$GITHUB_TOKEN" ]
then
  echo "Require environment variable: GITHUB_TOKEN"
  exit 1
fi

TAG=$(cat config/config.go | grep VERSION | awk -F'"' '{$0=$2}1')

git push --delete origin "v${TAG}"
git tag "v${TAG}"
git push --tags --force

curl --silent --output /dev/null --data "{\"tag_name\": \"v${TAG}\",\"target_commitish\": \"master\",\"name\": \"Pre Release v${TAG}\",\"body\": \"\",\"draft\": true,\"prerelease\": true}" https://api.github.com/repos/swapbyt3s/zenit/releases?access_token=${GITHUB_TOKEN}

ID=$(curl -sH "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/swapbyt3s/zenit/releases | jq -r '.[0].id')

rm -rf pkg/*

mkdir -p pkg/linux_amd64/
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o pkg/linux_amd64/zenit main.go
tar -czf pkg/linux_amd64/zenit-linux_amd64.tar.gz -C pkg/linux_amd64/ zenit

curl -# \
     --silent \
     --output /dev/null \
     -XPOST \
     -H "Authorization:token ${GITHUB_TOKEN}" \
     -H "Content-Type:application/octet-stream" \
     --data-binary @pkg/linux_amd64/zenit-linux_amd64.tar.gz \
     https://uploads.github.com/repos/swapbyt3s/zenit/releases/${ID}/assets?name=zenit-linux_amd64.tar.gz

echo -e "\r"
