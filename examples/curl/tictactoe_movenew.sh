#!/usr/bin/env bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

ROOM_ID="${1:?usage: $0 <room_id>}"

payload=$(cat <<EOF
{ "move" : [{
  "x": ["3"],
  "y": ["2"],
  "player": ["x"]
}]
}
EOF
)

curl_json POST "/room/$ROOM_ID/assert/move" "$payload" | jq .
