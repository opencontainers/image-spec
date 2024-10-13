#!/bin/sh

set -e
cd "$(dirname $0)/.."

if [ ! -x "$(command -v curl )" ] || [ ! -x "$(command -v jq )" ] ||  [ ! -x "$(command -v find )" ] ||  [ ! -x "$(command -v sed )" ]; then 
  echo "This command requires the following to run: curl, find, jq, and sed" >&2
  exit 1
fi

runtime_tag=$(curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/opencontainers/runtime-spec/releases/latest \
| jq -r .tag_name)

distribution_tag=$(curl -L \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/opencontainers/distribution-spec/releases/latest \
| jq -r .tag_name)

find . -name '*.md' -exec sed -i \
  -e "s#https://github.com/opencontainers/runtime-spec/blob/main/#https://github.com/opencontainers/runtime-spec/blob/${runtime_tag}/#g" \
  -e "s#https://github.com/opencontainers/distribution-spec/blob/main/#https://github.com/opencontainers/distribution-spec/blob/${distribution_tag}/#g" \
  '{}' \;
