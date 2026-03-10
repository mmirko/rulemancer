#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

payload=$(cat <<EOF
{
  "name": "example-brroom",
  "bridge_ref": "{{ .GameName }}"
}
EOF
)

curl_json POST "/brroom/create" "$payload" | jq .
